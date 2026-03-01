import { defineCollection, z } from 'astro:content';

// 博客文章集合
const blogCollection = defineCollection({
  type: 'content',
  schema: z.object({
    title: z.string(),
    description: z.string().optional(),
    pubDate: z.coerce.date(),
    updatedDate: z.coerce.date().optional(),
    heroImage: z.string().optional(),
    heroAlt: z.string().optional(),
    tags: z.array(z.string()).default([]),
    category: z.string().optional(),
    draft: z.boolean().default(false),
    author: z.string().default('Anonymous'),
    featured: z.boolean().default(false),
    // 用于评论区
    comments: z.boolean().default(true),
    // 文章唯一标识符，用于评论关联
    slug: z.string().optional(),
  }),
});

// 页面集合 (如关于页面等)
const pagesCollection = defineCollection({
  type: 'content',
  schema: z.object({
    title: z.string(),
    description: z.string().optional(),
    layout: z.string().optional(),
    showInNav: z.boolean().default(false),
    order: z.number().default(0),
  }),
});

export const collections = {
  blog: blogCollection,
  pages: pagesCollection,
};