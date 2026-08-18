import { getSecurityMetadata } from "@/utils.js";
import {schemaWithPagination, ZCluster, ZLink} from "@glimpse/zod";
import { initContract } from "@ts-rest/core";
import z from "zod";

const c = initContract();

const metadata = getSecurityMetadata();

export const clusterContract = c.router(
  {
    getClusters: {
      summary: "Get all clusters",
      path: "/clusters",
      method: "GET",
      description: "Get all clusters",
      query: z.object({
        page: z.number().min(1).optional(),
        limit: z.number().min(1).max(100).optional(),
        sort: z.enum(["created_at", "updated_at", "name"]).optional(),
        order: z.enum(["asc", "desc"]).optional(),
        search: z.string().min(1).optional(),
      }),
      responses: {
        200: schemaWithPagination(ZCluster)
      },
      metadata: metadata
    },

    getClusterById: {
    summary: "Get cluster by ID",
    path: "/clusters/:clusterId",
    method: "GET",
    description: "Get cluster by ID",
    responses: {
        200: ZCluster,
    },
    metadata: metadata,
    },

    createClusterLink: {
    summary: "Create a share link for a cluster",
    path: "/clusters/:clusterId/links",
    method: "POST",
    description: "Creates a shareable link for the given cluster. Supplying a password makes the link password protected.",
    body: z.object({
        password: z.string().min(6).optional(),
        expiresAt: z.string().datetime().nullish(),
    }),
    responses: {
        201: ZLink,
    },
    metadata: metadata,
    }
  },
  {
    pathPrefix: "/api/v1"
  }
)
