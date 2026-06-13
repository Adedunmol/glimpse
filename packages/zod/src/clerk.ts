import z from "zod";


export const ZClerkEventPayload = z.object({
  data: z.object({}).passthrough(),
  object: z.string(),
  type: z.string(),
  timestamp: z.string(),
  instance_id: z.string(),
});