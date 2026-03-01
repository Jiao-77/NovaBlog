// @ts-check
import { defineConfig } from 'astro/config';
import mdx from '@astrojs/mdx';
import vue from '@astrojs/vue';
import tailwind from '@astrojs/tailwind';
import react from '@astrojs/react';
import remarkMath from 'remark-math';
import rehypeKatex from 'rehype-katex';

// https://astro.build/config
export default defineConfig({
  site: 'https://your-domain.com', // 替换为实际域名
  integrations: [
    mdx({
      // 配置 MDX 组件映射
      optimize: true,
    }),
    vue(),
    react(),
    tailwind({
      applyBaseStyles: false, // 我们将手动控制基础样式
    }),
  ],
  markdown: {
    shikiConfig: {
      theme: 'github-dark',
      wrap: true,
    },
    remarkPlugins: [remarkMath],
    rehypePlugins: [rehypeKatex],
  },
  vite: {
    ssr: {
      noExternal: ['@vueuse/core'],
    },
  },
  // 配置构建输出
  output: 'static',
  build: {
    assets: 'assets',
  },
});
