-- สร้างตาราง courses
CREATE TABLE courses (
    course_id SERIAL PRIMARY KEY,
    course_name VARCHAR(255) NOT NULL,
    description TEXT,
    thumbnail_url VARCHAR(255),
    instructor_name VARCHAR(255),
    profile_url VARCHAR(255),
    duration VARCHAR(255),
    price DECIMAL(10,2),
    detail_url VARCHAR(255),
    rating DECIMAL(2,1),
    num_reviews INT,
    enrollment_count INT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- สร้าง function สำหรับอัพเดท created_at โดยอัตโนมัติ
CREATE OR REPLACE FUNCTION created_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- สร้าง function สำหรับอัพเดท updated_at โดยอัตโนมัติ
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- สร้าง trigger สำหรับอัพเดท updated_at
CREATE TRIGGER update_courses_modtime
    BEFORE UPDATE ON courses
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();


CREATE TABLE affiliates (
    affiliate_id SERIAL PRIMARY KEY,
    affiliate_name VARCHAR(255),
    affiliate_email VARCHAR(255) UNIQUE,
    affiliate_password VARCHAR(255),
    affiliate_url VARCHAR(255),
    affiliate_api_key VARCHAR(255) UNIQUE
)

CREATE TABLE request_logs (
    id SERIAL PRIMARY KEY,
    affiliate_id VARCHAR(255),
    action VARCHAR(255),
    parameter TEXT,
    timestamp TIMESTAMP
)

CREATE TABLE Click_logs (
    affiliate_id SERIAL PRIMARY KEY,
    course_id SERIAL ,
    click_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    clicks INTEGER,
    FOREIGN KEY (affiliate_id) REFERENCES affiliates(id),
    FOREIGN KEY (course_id) REFERENCES courses(course_id)
)

-- Course_ID         string  `json:"-"`
-- 	Course_Name       string  `json:"-"`
-- 	Course_Desc       string  `json:"-"` // คำอธิบาย
-- 	detail_url        string  `json:"-"`
-- 	thumbnail_url     string  `json:"-"`
-- 	Course_Type       string  `json:"-"`
-- 	Course_Instructor string  `json:"-"`
-- 	profile_url       string  `json:"-"`
-- 	Course_Price      int     `json:"-"`
-- 	duration          int     `json:"-"`
-- 	rating            float64 `json:"-"`
-- 	num_reviews       int     `json:"-"`
-- 	enrollment_count  int     `json:"-"`
-- 	created_at        string  `json:"-"`
-- 	updated_at        string  `json:"-"`

-- เพิ่มข้อมูลตัวอย่าง
INSERT INTO courses (
    course_name, -- ชื่อคอร์ส
    course_Desc, -- คำอธิบายเกี่ยวกับคอร์ส
    detail_url, -- ลิงก์ไปยังหน้ารายละเอียดของคอร์ส
    thumbnail_url, -- URL ของภาพปกหรือรูปตัวอย่างของคอร์ส
    course_type, --ประเภทของคอร์ส 
    course_instructor, --ชื่อของผู้สอน
    profile_url, -- ลิงก์ไปยังโปรไฟล์ของผู้สอน
    course_price, -- ราคาของคอร์ส
    duration, -- ระยะเวลาเรียน
    rating, -- คะแนนรีวิวเฉลี่ย
    num_reviews, -- จำนวนรีวิว
    enrollment_count -- จำนวนผู้ลงทะเบียนเรียนแล้ว
) VALUES (
    'Python for Beginners',
    'เรียนรู้พื้นฐานการเขียนโปรแกรมด้วยภาษา Python เหมาะสำหรับผู้เริ่มต้น',
    'https://example.com/python-beginners',
    'https://example.com/images/python.jpg',
    'Python',
    'Jane Smith',
    'https://example.com/instructors/jane',
    0.00,
    '6 ชั่วโมง',
    4.8,
    1523,
    28000
),

('Excel ขั้นเทพสำหรับงานออฟฟิศ', 
 'ใช้งาน Excel ตั้งแต่พื้นฐานถึงขั้นสูง พร้อมสูตรและเทคนิค', 
 'https://example.com/excel-master', 
 'https://example.com/images/excel.jpg',
 'Excel',
 'ภัทรพล โพธิ์ศรี', 
 'https://example.com/instructors/pat', 
 699.00, 
 '8 ชั่วโมง', 
 4.6, 
 875, 
 12500),

('วาดภาพ Digital Art ด้วย Procreate',
 'เรียนรู้การวาดภาพบน iPad ด้วยแอป Procreate ตั้งแต่พื้นฐาน',
 'https://example.com/procreate-course',
 'https://example.com/images/procreate.jpg',
 'Art',
 'Arisa Chen',
 'https://example.com/instructors/arisa',
 1200.00,
 '5 ชั่วโมง 30 นาที',
 4.9,
 412,
 5300
),

('JavaScript Web Development',
 'เขียนเว็บไซต์แบบ Interactive ด้วย JavaScript',
 'https://example.com/js-webdev',
 'https://example.com/images/js.jpg',
 'JavaScript',
 'John Doe',
 'https://example.com/instructors/john',
 0.00,
 '7 ชั่วโมง',
 4.7,
 1399,
 24000
),

('การจัดการเวลาอย่างมีประสิทธิภาพ',
 'เรียนรู้เทคนิค Time Management และ Productivity',
 'https://example.com/time-management',
 'https://example.com/images/time.jpg',
 'Management',
 'นภัสสร ลิ้มวัฒนา',
 'https://example.com/instructors/napas0',
 499.00,
 '3 ชั่วโมง',
 4.5,
 326,
 8000
),

('UI/UX Design สำหรับมือใหม่',
 'เรียนรู้การออกแบบประสบการณ์ผู้ใช้ในแอปและเว็บไซต์',
 'https://example.com/uiux-design',
 'https://example.com/images/uiux.jpg',
 'UX/UI',
 'Kelvin Wong',
 'https://example.com/instructors/kelvin',
 990.00,
 '10 ชั่วโมง',
 4.8,
 978,
 10800
),

('ภาษาอังกฤษสำหรับการทำงาน',
 'พัฒนาทักษะภาษาอังกฤษในที่ทำงาน เน้นการสนทนา',
 'https://example.com/business-english',
 'https://example.com/images/english.jpg',
 'English',
 'ครูแพร',
 'https://example.com/instructors/kru-prae',
 499.00,
 '4 ชั่วโมง',
 4.4,
 610,
 19400
),

('เรียนถ่ายภาพด้วยกล้องมือถือ',
),
;

 INSERT INTO affiliates(
    affiliate_id,
    affiliate_name,
    affiliate_email,
    affiliate_password,
    affiliate_url,
    affiliate_api_key
)VALUES (
    'affiliate123',
    'affiliate123@example.com',
    'password123',
    'https://example.com/affiliate123',
    'jhksfgolzsdflkhsdfs665sgghdzjrksbdfzxhb2145fsgvsoe8twuhfuijdcvoiaErfqo3ue' -- พิมพ์มั่วนะ
);