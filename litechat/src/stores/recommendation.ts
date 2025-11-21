import { defineStore } from 'pinia';
import { ref } from 'vue';
import {
    getKnowledgeBaseList,
    getQuestionList,
    getRecommendationConfig,
    type KnowledgeBaseItem,
    type QuestionItem
} from '../api/recommendation';

export const useRecommendationStore = defineStore('recommendation', () => {
    const knowledgeBaseList = ref<KnowledgeBaseItem[]>([]);
    const questionList = ref<QuestionItem[]>([]);
    const showKnowledgeBase = ref(true);
    const showQuestions = ref(true);
    const isLoading = ref(false);

    const fetchRecommendations = async () => {
        isLoading.value = true;
        try {
            const [kbList, qList, config] = await Promise.all([
                getKnowledgeBaseList(),
                getQuestionList(),
                getRecommendationConfig()
            ]);

            knowledgeBaseList.value = kbList;
            questionList.value = qList;
            showKnowledgeBase.value = config.showKnowledgeBase;
            showQuestions.value = config.showQuestions;
        } catch (error) {
            console.error('Failed to fetch recommendations:', error);
        } finally {
            isLoading.value = false;
        }
    };

    // Helper to get random items
    const getRandomItems = <T>(list: T[], count: number): T[] => {
        const shuffled = [...list].sort(() => 0.5 - Math.random());
        return shuffled.slice(0, count);
    };

    return {
        knowledgeBaseList,
        questionList,
        showKnowledgeBase,
        showQuestions,
        isLoading,
        fetchRecommendations,
        getRandomItems
    };
});
