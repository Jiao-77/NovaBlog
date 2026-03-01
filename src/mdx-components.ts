import Counter from './components/Counter.vue';
import TypstBlock from './components/TypstBlock.astro';

/**
 * MDX 组件映射
 * 
 * 在 MDX 文件中可以直接使用这些组件，无需 import。
 * 例如：
 * ```mdx
 * <Counter client:visible />
 * <TypstBlock>$ E = mc^2 $</TypstBlock>
 * ```
 */
export const components = {
  // Vue 交互组件
  Counter,
  
  // Astro 静态组件
  TypstBlock,
};