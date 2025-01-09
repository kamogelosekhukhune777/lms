-- Version: 1.01
CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    user_name TEXT NOT NULL,
    user_email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL
);

-- Version: 1.02
CREATE TABLE lectures (
    lecture_id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    video_url TEXT NOT NULL,
    public_id TEXT NOT NULL,
    free_preview BOOLEAN NOT NULL DEFAULT FALSE
);

-- Version: 1.03
CREATE TABLE courses (
    course_id UUID PRIMARY KEY,
    instructor_id UUID NOT NULL,
    title TEXT NOT NULL,
    category TEXT NOT NULL,
    level TEXT NOT NULL,
    primary_language TEXT NOT NULL,
    subtitle TEXT,
    description TEXT NOT NULL,
    image TEXT NOT NULL,
    welcome_message TEXT,
    pricing NUMERIC(10, 2) NOT NULL,
    objectives TEXT,
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    date TIMESTAMP NOT NULL,
    FOREIGN KEY (instructor_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Version: 1.04
CREATE TABLE curriculum (
    curriculum_id UUID PRIMARY KEY,
    course_id UUID NOT NULL,
    lecture_id UUID NOT NULL,
    FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE,
    FOREIGN KEY (lecture_id) REFERENCES lectures(lecture_id) ON DELETE CASCADE
);

-- Version: 1.05
CREATE TABLE course_students (
    course_student_id UUID PRIMARY KEY,
    course_id UUID NOT NULL,
    student_id UUID NOT NULL,
    paid_amount NUMERIC(10, 2) NOT NULL,
    FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE,
    FOREIGN KEY (student_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Version: 1.06
CREATE TABLE student_courses (
    student_course_id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    course_id UUID NOT NULL,
    date_of_purchase TIMESTAMP NOT NULL,
    course_image TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE
);

-- Version: 1.07
CREATE TABLE orders (
    order_id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    order_status TEXT NOT NULL,
    payment_method TEXT NOT NULL,
    payment_status TEXT NOT NULL,
    order_date TIMESTAMP NOT NULL,
    payment_id TEXT,
    payer_id TEXT,
    course_id UUID NOT NULL,
    course_pricing NUMERIC(10, 2) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE
);

-- Version: 1.08
CREATE TABLE lectures_progress (
    lecture_progress_id UUID PRIMARY KEY,
    lecture_id UUID NOT NULL,
    viewed BOOLEAN NOT NULL DEFAULT FALSE,
    date_viewed TIMESTAMP,
    FOREIGN KEY (lecture_id) REFERENCES lectures(lecture_id) ON DELETE CASCADE
);

-- Version: 1.09
CREATE TABLE course_progress (
    course_progress_id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    course_id UUID NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    completion_date TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE
);
