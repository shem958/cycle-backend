package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shem958/cycle-backend/config"
	"github.com/shem958/cycle-backend/models"
)

// CreatePost creates a new post by the authenticated user
func CreatePost(c *gin.Context) {
	var input struct {
		Title       string   `json:"title"`
		Content     string   `json:"content"`
		Tags        []string `json:"tags"`
		IsAnonymous bool     `json:"is_anonymous"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	post := models.Post{
		ID:          uuid.New(),
		AuthorID:    userID.(uuid.UUID),
		Title:       input.Title,
		Content:     input.Content,
		Tags:        input.Tags,
		IsAnonymous: input.IsAnonymous,
	}

	if err := config.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// GetAllPosts retrieves posts with optional filtering, sorting, and searching
func GetAllPosts(c *gin.Context) {
	tagFilter := c.Query("tag")
	search := c.Query("search")
	sort := c.DefaultQuery("sort", "new") // "new" or "top"

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID)

	blockedIDs, err := getBlockedUserIDs(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get blocked users"})
		return
	}

	db := config.DB.Model(&models.Post{}).Preload("Comments")

	// Filter out posts from blocked users
	if len(blockedIDs) > 0 {
		db = db.Where("author_id NOT IN ?", blockedIDs)
	}

	if tagFilter != "" {
		db = db.Where("? = ANY (tags)", tagFilter)
	}

	if search != "" {
		pattern := "%" + search + "%"
		db = db.Where("title ILIKE ? OR content ILIKE ?", pattern, pattern)
	}

	switch sort {
	case "top":
		db = db.Order("array_length(comments, 1) DESC")
	default:
		db = db.Order("created_at DESC")
	}

	var posts []models.Post
	if err := db.Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// GetPostByID retrieves a single post by ID
func GetPostByID(c *gin.Context) {
	postID := c.Param("id")

	var post models.Post
	if err := config.DB.Preload("Comments.Replies").First(&post, "id = ?", postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	userID, _ := c.Get("user_id")
	userUUID, _ := userID.(uuid.UUID)

	// Add post reaction stats
	var postReactions []struct {
		Type  string
		Count int
	}
	config.DB.Model(&models.Reaction{}).
		Select("type, COUNT(*) as count").
		Where("target_id = ? AND target_type = ?", post.ID, "post").
		Group("type").
		Scan(&postReactions)

	for _, r := range postReactions {
		if r.Type == "like" {
			post.LikeCount = r.Count
		} else if r.Type == "dislike" {
			post.DislikeCount = r.Count
		}
	}

	// Check if current user reacted
	var userPostReaction models.Reaction
	if err := config.DB.Where("user_id = ? AND target_id = ? AND target_type = ?", userUUID, post.ID, "post").First(&userPostReaction).Error; err == nil {
		post.UserReaction = userPostReaction.Type
	}

	// Comments and replies reactions
	for i := range post.Comments {
		comment := &post.Comments[i]

		// Count reactions
		var counts []struct {
			Type  string
			Count int
		}
		config.DB.Model(&models.Reaction{}).
			Select("type, COUNT(*) as count").
			Where("target_id = ? AND target_type = ?", comment.ID, "comment").
			Group("type").
			Scan(&counts)

		for _, r := range counts {
			if r.Type == "like" {
				comment.LikeCount = r.Count
			} else if r.Type == "dislike" {
				comment.DislikeCount = r.Count
			}
		}

		// User's reaction
		var userCommentReaction models.Reaction
		if err := config.DB.Where("user_id = ? AND target_id = ? AND target_type = ?", userUUID, comment.ID, "comment").First(&userCommentReaction).Error; err == nil {
			comment.UserReaction = userCommentReaction.Type
		}

		// Repeat for replies
		for j := range comment.Replies {
			reply := &comment.Replies[j]
			var replyCounts []struct {
				Type  string
				Count int
			}
			config.DB.Model(&models.Reaction{}).
				Select("type, COUNT(*) as count").
				Where("target_id = ? AND target_type = ?", reply.ID, "comment").
				Group("type").
				Scan(&replyCounts)

			for _, r := range replyCounts {
				if r.Type == "like" {
					reply.LikeCount = r.Count
				} else if r.Type == "dislike" {
					reply.DislikeCount = r.Count
				}
			}

			var userReplyReaction models.Reaction
			if err := config.DB.Where("user_id = ? AND target_id = ? AND target_type = ?", userUUID, reply.ID, "comment").First(&userReplyReaction).Error; err == nil {
				reply.UserReaction = userReplyReaction.Type
			}
		}
	}

	c.JSON(http.StatusOK, post)
}

// GetAllTags returns a list of unique tags from all posts
func GetAllTags(c *gin.Context) {
	var tags []string
	db := config.DB

	if err := db.
		Raw(`SELECT DISTINCT unnest(tags) FROM posts`).
		Scan(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
		return
	}

	c.JSON(http.StatusOK, tags)
}

// CreateComment adds a comment to a post
func CreateComment(c *gin.Context) {
	var input struct {
		PostID      uuid.UUID `json:"post_id"`
		Content     string    `json:"content"`
		IsAnonymous bool      `json:"is_anonymous"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	comment := models.Comment{
		ID:          uuid.New(),
		PostID:      input.PostID,
		AuthorID:    userID.(uuid.UUID),
		Content:     input.Content,
		IsAnonymous: input.IsAnonymous,
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// ReplyToComment handles replying to a comment (nested)
func ReplyToComment(c *gin.Context) {
	var input struct {
		PostID      uuid.UUID `json:"post_id" binding:"required"`
		ParentID    uuid.UUID `json:"parent_id" binding:"required"`
		Content     string    `json:"content" binding:"required"`
		IsAnonymous bool      `json:"is_anonymous"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	reply := models.Comment{
		ID:          uuid.New(),
		PostID:      input.PostID,
		AuthorID:    userID.(uuid.UUID),
		Content:     input.Content,
		IsAnonymous: input.IsAnonymous,
		ParentID:    &input.ParentID,
	}

	if err := config.DB.Create(&reply).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reply"})
		return
	}

	c.JSON(http.StatusCreated, reply)
}

// ReportContent allows a user to report a post or comment
func ReportContent(c *gin.Context) {
	var input struct {
		TargetPostID    *uuid.UUID `json:"target_post_id"`
		TargetCommentID *uuid.UUID `json:"target_comment_id"`
		Reason          string     `json:"reason"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.TargetPostID == nil && input.TargetCommentID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Target post or comment must be specified"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	report := models.Report{
		ID:              uuid.New(),
		ReporterID:      userID.(uuid.UUID),
		TargetPostID:    input.TargetPostID,
		TargetCommentID: input.TargetCommentID,
		Reason:          input.Reason,
	}

	if err := config.DB.Create(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to report content"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Content reported successfully"})
}

func getBlockedUserIDs(userID uuid.UUID) ([]uuid.UUID, error) {
	var blocked []models.Block
	if err := config.DB.Where("user_id = ?", userID).Find(&blocked).Error; err != nil {
		return nil, err
	}

	var blockedIDs []uuid.UUID
	for _, b := range blocked {
		blockedIDs = append(blockedIDs, b.TargetID)
	}
	return blockedIDs, nil
}
