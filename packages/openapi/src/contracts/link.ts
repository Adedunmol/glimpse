import { getSecurityMetadata } from "@/utils.js";
import {schemaWithPagination, ZLink} from "@glimpse/zod";
import { initContract } from "@ts-rest/core";
import z from "zod";

const c = initContract();

const metadata = getSecurityMetadata();

const paginationQuery = z.object({
  page: z.number().min(1).optional(),
  limit: z.number().min(1).max(100).optional(),
  sort: z.enum(["created_at", "updated_at", "name"]).optional(),
  order: z.enum(["asc", "desc"]).optional(),
  search: z.string().min(1).optional(),
});

export const linkContract = c.router(
  {
    getLinks: {
      summary: "Get all links",
      path: "/links",
      method: "GET",
      description: "Get all links",
      query: paginationQuery,
      responses: {
        200: schemaWithPagination(ZLink)
      },
      metadata: metadata
    },

    getLinkById: {
    summary: "Get link by ID",
    path: "/links/id/:linkId",
    method: "GET",
    description: "Get a single link by its id",
    responses: {
        200: ZLink,
    },
    metadata: metadata,
    },

    getLinkByToken: {
    summary: "Get link by token",
    path: "/links/token/:token",
    method: "GET",
    description: "Get a single link by its share token",
    responses: {
        200: ZLink,
    },
    metadata: metadata,
    },

    getLinksByClusterId: {
    summary: "Get links by cluster ID",
    path: "/links/:clusterId",
    method: "GET",
    description: "Get the links belonging to a cluster",
    query: paginationQuery,
    responses: {
        200: schemaWithPagination(ZLink),
    },
    metadata: metadata,
    }
  },
  {
    pathPrefix: "/api/v1"
  }
)
