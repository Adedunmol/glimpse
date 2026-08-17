import z from "zod";


export const ZClerkEventPayload = z.object({
  data: z.object({}).passthrough(),
  object: z.string(),
  type: z.string(),
  timestamp: z.number(), // Unix milliseconds
  instance_id: z.string().optional(),
});