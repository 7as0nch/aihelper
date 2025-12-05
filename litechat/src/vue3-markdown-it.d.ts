declare module 'vue3-markdown-it' {
    import { DefineComponent } from 'vue';
    const Vue3MarkdownIt: DefineComponent<{
        source: string;
        plugins?: { plugin: any; options?: any }[];
        html?: boolean;
        linkify?: boolean;
        typographer?: boolean;
        highlight?: ((str: string, lang: string) => string) | { highlighter: (str: string, lang: string) => string };
    }>;
    export default Vue3MarkdownIt;
}
