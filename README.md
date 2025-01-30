• Developed a recruitment system utilizing GORM for ORM, JWT for secure authentication, and role-based access control, ensuring seamless and secure user management.
• Designed a normalized database structure with six interlinked tables—user, job, resume, education, experience, and job application—allowing for efficient data retrieval and management.
• Built robust APIs for key functionalities, including user signup and login, job posting and application management, and fetching details for applicants and job postings.
• Integrated third-party resume parsing with the UploadResume endpoint, enabling automatic data extraction from uploaded resumes.
• Implemented middleware for JWT verification and employer authorization, enhancing security and enforcing access control.
• Key APIs:

Signup, Login: User authentication and registration.
GetApplicantData, GetAllApplicants, GetAllResumes: Fetch applicant profiles and resumes.
GetMyJobsDetail, AddJob, GetAllJobs, GetJobData: Manage job postings and retrieve job information.
ApplyToJob: Allow applicants to submit job applications.

All the schema relationships are depicted down below: 
<img width="859" alt="schemas" src="https://github.com/user-attachments/assets/80d21afe-3a12-4013-9b4b-94ea667b5881">
