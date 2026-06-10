import z from "zod";

export const ZUploadStatus = z.enum(["pending", "processing", "done", "failed"]);

export const ZUpload = z.object({
    id: z.string().uuid(),
    name: z.string(),
    hostId: z.string(),
    status: ZUploadStatus,
    expiresAt: z.string(),
    createdAt: z.string(),
    updatedAt: z.string(),
})