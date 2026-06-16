import { getSecurityMetadata } from "@/utils.js";
import {schemaWithPagination, ZLink} from "@glimpse/zod";
import { initContract } from "@ts-rest/core";
import z from "zod";

const c = initContract();

const metadata = getSecurityMetadata();

export const linkContract = c.router(
  {
    getLinks: {
      summary: "Get all links",
      path: "/links",
      method: "GET",
      description: "Get all links",
      query: z.object({
        page: z.number().min(1).optional(),
        limit: z.number().min(1).max(100).optional(),
        sort: z.enum(["created_at", "updated_at", "name"]).optional(),
        order: z.enum(["asc", "desc"]).optional(),
        search: z.string().min(1).optional(),
      }),
      responses: {
        200: schemaWithPagination(ZLink)
      },
      metadata: metadata
    },

    getLinkById: {
    summary: "Get link by ID",
    path: "/clusters/:linkId",
    method: "GET",
    description: "Get link by ID",
    responses: {
        200: ZLink,
    },
    metadata: metadata,
    },

    getLinkByClusterId: {
    summary: "Get link by cluster ID",
    path: "/clusters/:clusterId",
    method: "GET",
    description: "Get link by cluster ID",
    responses: {
        200: ZLink,
    },
    metadata: metadata,
    },

    getLinkByToken: {
    summary: "Get link by token",
    path: "/clusters/:token",
    method: "GET",
    description: "Get link by token",
    responses: {
        200: ZLink,
    },
    metadata: metadata,
    }
  }
)