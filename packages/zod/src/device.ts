import z from "zod"

export const ZUserDevice = z.object({
    id: z.string().uuid(),
    userId: z.string(),
    pushToken: z.string(),
    platform: z.string(),
    expiresAt: z.string(),
    createdAt: z.string(),
    updatedAt: z.string(),
})