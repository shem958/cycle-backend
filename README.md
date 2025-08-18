# 🌸 Cycle Tracking Backend

This is the **backend service** for a Cycle Tracking & Women’s Health App.  
It provides APIs for managing menstrual cycles, pregnancy & postpartum monitoring, medical appointments, and community features — with a strong focus on **data privacy and encryption**.

---

## 🚀 Features (Implemented So Far)

- **Authentication & Authorization**
  - User registration & login
  - JWT-based authentication
  - Admin middleware & role-based access

- **User Management**
  - Profile creation & updates
  - Blocking, muting, suspending users (admin tools)

- **Cycle Tracking**
  - Add & manage menstrual cycle records
  - Track mood & symptoms
  - Predictive cycle insights:
    - Average cycle length
    - Next period prediction
    - Ovulation & fertile window
    - Mood & symptom patterns

- **Pregnancy & Postpartum Monitoring**
  - Record pregnancy & postpartum health data
  - Encrypted storage of sensitive medical info
  - Secure retrieval for the logged-in user only

- **Medical Appointments / Follow-Up Scheduling**
  - Create & manage doctor appointments
  - Store appointment notes and reminders

- **Community Features**
  - Posts, comments, and reporting
  - Admin moderation of posts/comments
  - Report management dashboard

---

## 📌 Upcoming Features

- 🔲 **Notifications & Reminders** (appointments, cycle phase alerts)  
- 🔲 **Export Health Data** (PDF/CSV reports)  
- 🔲 **Doctor–Patient Messaging** (secure encrypted chat)  
- 🔲 **Multi-language Support**  
- 🔲 **Accessibility Enhancements** (screen readers, high contrast)  
- 🔲 **In-App Support / Feedback System**

---

## 🛠️ Tech Stack

- **Language:** Go (Golang)  
- **Framework:** [Gin](https://github.com/gin-gonic/gin)  
- **Database:** PostgreSQL (via GORM ORM)  
- **Authentication:** JWT-based  
- **Security:** AES-256 data encryption for sensitive health info  
- **Middleware:** CORS, Authentication, Role-based Access  

---

## ⚙️ Installation & Setup

### 1️⃣ Install Go (if missing or after OS update)

```bash
# Remove old Go version (if any)
sudo rm -rf /usr/local/go

# Download & install Go 1.22+ (latest stable)
wget https://go.dev/dl/go1.22.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz

# Add Go to PATH
echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
source ~/.bashrc

# Verify installation
go version
