import { initContract } from "@ts-rest/core";
import { z } from "zod";
import { ZClerkEventPayload } from "@glimpse/zod";

const c = initContract();

export const clerkWebHookContract = c.router({
  postEvent: {
    summary: "Handle Clerk webhook events",
    path: "/clerk/webhook",
    method: "POST",
    description:
      "Verifies Clerk webhook signatures and synchronizes local user data.",
    body: ZClerkEventPayload,
    responses: {
      200: z.void(),
    },
  },
});