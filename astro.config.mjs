// @ts-check
import { defineConfig } from 'astro/config';
import mdx from '@astrojs/mdx';
import vue from '@astrojs/vue';
import tailwind from '@astrojs/tailwind';

// https://astro.build/config
export default defineConfig({
  site: 'https://your-domain.com', // 替换为实际域名
  integrations: [
    mdx({
      // 配置 MDX 组件映射
      optimize: true,
    }),
    vue(),
    tailwind({
      applyBaseStyles: false, // 我们将手动控制基础样式
    }),
  ],
  markdown: {
    shikiConfig: {
      theme: 'github-dark',
      wrap: true,
    },
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
