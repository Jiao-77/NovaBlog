/// <reference path="../.astro/types.d.ts" />
/// <reference types="astro/client" />

// Vue 组件类型声明
declare module '*.vue' {
  import type { DefineComponent } from 'vue';
  const component: DefineComponent<{}, {}, any>;
  export default component;
}

// MDX 文件类型声明
declare module '*.mdx' {
  import type { MDXComponent } from 'mdx/types';
  export default MDXComponent;
}