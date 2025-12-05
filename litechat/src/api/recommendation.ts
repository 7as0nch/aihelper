
export interface KnowledgeBaseItem {
    id: string;
    title: string;
    description: string;
    icon: string; // URL or icon name
}

export interface QuestionItem {
    id: string;
    content: string;
}

export interface RecommendationConfig {
    showKnowledgeBase: boolean;
    showQuestions: boolean;
}

// Mock Data
const mockKnowledgeBase: KnowledgeBaseItem[] = [
    { id: '1', title: '正经的大学生活与不正经的PC', description: '关于大学生活和电脑知识的探讨', icon: 'BookOpen' },
    { id: '2', title: '具身智能 VLA+RL', description: '前沿 AI 技术研究', icon: 'Brain' },
    { id: '3', title: '小K配方师护肤知识库', description: '科学护肤指南', icon: 'FlaskConical' },
    { id: '4', title: '前端工程化实践', description: '现代前端开发工作流', icon: 'Code' },
    { id: '5', title: '摄影构图技巧', description: '提升摄影审美的关键', icon: 'Camera' },
];

const mockQuestions: QuestionItem[] = [
    { id: '1', content: '房价下跌对普通人意味着什么？' },
    { id: '2', content: '性价比高的显卡有哪些？' },
    { id: '3', content: '改善睡眠质量的具体方法有哪些？' },
    { id: '4', content: '以「低空经济」为主题构建思论文框架' },
    { id: '5', content: '如何评价 Vue 3 的 Composition API？' },
    { id: '6', content: '2025 年有哪些值得关注的科技趋势？' },
    { id: '7', content: '如何快速入门 Python 数据分析？' },
];

export const getKnowledgeBaseList = async (): Promise<KnowledgeBaseItem[]> => {
    await new Promise(resolve => setTimeout(resolve, 500));
    return mockKnowledgeBase;
};

export const getQuestionList = async (): Promise<QuestionItem[]> => {
    await new Promise(resolve => setTimeout(resolve, 500));
    return mockQuestions;
};

export const getRecommendationConfig = async (): Promise<RecommendationConfig> => {
    await new Promise(resolve => setTimeout(resolve, 300));
    return {
        showKnowledgeBase: false,
        showQuestions: true,
    };
};
