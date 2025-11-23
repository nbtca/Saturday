-- Create member_applications table for tracking new member applications
CREATE TABLE IF NOT EXISTS public.member_applications (
    application_id VARCHAR(36) PRIMARY KEY,
    member_id VARCHAR(10) NOT NULL,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    section VARCHAR(50) NOT NULL,
    qq VARCHAR(20),
    email VARCHAR(100),
    major VARCHAR(100),
    class VARCHAR(50),
    memo TEXT,

    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),

    reviewed_by VARCHAR(10),
    reviewed_at TIMESTAMP WITHOUT TIME ZONE,
    reject_reason TEXT,

    gmt_create TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    gmt_modified TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

-- Add indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_member_applications_status ON public.member_applications(status);
CREATE INDEX IF NOT EXISTS idx_member_applications_member_id ON public.member_applications(member_id);
CREATE INDEX IF NOT EXISTS idx_member_applications_gmt_create ON public.member_applications(gmt_create);

-- Add comments for documentation
COMMENT ON TABLE public.member_applications IS '成员申请表 (Member application records)';
COMMENT ON COLUMN public.member_applications.application_id IS '申请ID (Application ID)';
COMMENT ON COLUMN public.member_applications.member_id IS '学号 (Student ID)';
COMMENT ON COLUMN public.member_applications.name IS '姓名 (Full name)';
COMMENT ON COLUMN public.member_applications.phone IS '手机号 (Phone number)';
COMMENT ON COLUMN public.member_applications.section IS '部门 (Department/Section)';
COMMENT ON COLUMN public.member_applications.qq IS 'QQ号 (QQ number)';
COMMENT ON COLUMN public.member_applications.email IS '邮箱 (Email address)';
COMMENT ON COLUMN public.member_applications.major IS '专业 (Major)';
COMMENT ON COLUMN public.member_applications.class IS '班级 (Class)';
COMMENT ON COLUMN public.member_applications.memo IS '备注/自我介绍 (Memo/Self introduction)';
COMMENT ON COLUMN public.member_applications.status IS '状态: pending(待审核), approved(已批准), rejected(已拒绝)';
COMMENT ON COLUMN public.member_applications.reviewed_by IS '审核人 (Reviewed by member ID)';
COMMENT ON COLUMN public.member_applications.reviewed_at IS '审核时间 (Review timestamp)';
COMMENT ON COLUMN public.member_applications.reject_reason IS '拒绝原因 (Rejection reason)';
