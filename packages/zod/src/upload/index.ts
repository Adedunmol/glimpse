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
    photos: z.array(z.object({
        id: z.string().uuid(),
        uploadId: z.string().uuid(),
        storageKey: z.string(),
        status: z.enum(["pending", "uploaded"]),
        isEmbedded: z.boolean(),
        createdAt: z.string(),
        updatedAt: z.string(),
    }))
})

export const ZPresignedUrls =  z.object({
    uploadId: z.string(),
    uploads: z.array(z.object({
        key: z.string(),
        url: z.string().url()
    }))
})