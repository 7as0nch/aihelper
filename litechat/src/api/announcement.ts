export interface Announcement {
    id: string;
    content: string;
    type: 'info' | 'warning' | 'success';
}

export const fetchAnnouncement = async (): Promise<Announcement | null> => {
    // Simulate API delay
    await new Promise(resolve => setTimeout(resolve, 500));

    // Return a mock announcement
    return {
        id: 'announcement-2',
        content: '欢迎使用 AI 智能助手！我们刚刚更新了新的模型，支持更强大的代码生成能力。',
        type: 'info'
    };
};
