-- Version: 1.01
-- Description: Create table users
CREATE TABLE Users (
    user_id UUID PRIMARY KEY,
    user_name VARCHAR(255) NOT NULL,
    user_email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Version: 1.02
-- Courses Table
CREATE TABLE Courses (
    course_id UUID PRIMARY KEY,
    instructor_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    category VARCHAR(100),
    level VARCHAR(50),
    primary_language VARCHAR(50),
    subtitle TEXT,
    description TEXT,
    image TEXT,
    welcome_message TEXT,
    pricing DECIMAL(10,2) NOT NULL,
    objectives TEXT,
    is_published BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (instructor_id) REFERENCES Users(user_id) ON DELETE CASCADE
);

-- Version: 1.03
-- Lectures Table
CREATE TABLE Lectures (
    lecture_id UUID PRIMARY KEY,
    course_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    video_url TEXT NOT NULL,
    public_id VARCHAR(255) UNIQUE,
    free_preview BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (course_id) REFERENCES Courses(course_id) ON DELETE CASCADE
);

-- Version: 1.04
-- Enrollments Table (Students enrolled in courses)
CREATE TABLE Enrollments (
    enrollment_id UUID PRIMARY KEY,
    student_id UUID NOT NULL,
    course_id UUID NOT NULL,
    paid_amount DECIMAL(10,2) NOT NULL,
    enrolled_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (student_id) REFERENCES Users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES Courses(course_id) ON DELETE CASCADE
);

-- Version: 1.05
-- Orders Table
CREATE TABLE Orders (
    order_id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    order_status VARCHAR(50) NOT NULL CHECK (order_status IN ('pending', 'completed', 'failed', 'refunded')),
    payment_method VARCHAR(50) NOT NULL CHECK (payment_method IN ('credit_card', 'paypal', 'stripe', 'bank_transfer')),
    payment_status VARCHAR(50) NOT NULL CHECK (payment_status IN ('paid', 'unpaid', 'refunded')),
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    payment_id VARCHAR(255),
    payer_id VARCHAR(255),
    instructor_id UUID NOT NULL,
    course_id UUID NOT NULL,
    course_pricing DECIMAL(10,2) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES Users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (instructor_id) REFERENCES Users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES Courses(course_id) ON DELETE CASCADE
);

-- Version: 1.06
-- Course Progress Table
CREATE TABLE CourseProgress (
    progress_id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    course_id UUID NOT NULL,
    completed BOOLEAN DEFAULT FALSE,
    completion_date TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES Courses(course_id) ON DELETE CASCADE
);

-- Version: 1.07
-- Lecture Progress Table
CREATE TABLE LectureProgress (
    progress_id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    lecture_id UUID NOT NULL,
    viewed BOOLEAN DEFAULT FALSE,
    date_viewed TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (lecture_id) REFERENCES Lectures(lecture_id) ON DELETE CASCADE
);
