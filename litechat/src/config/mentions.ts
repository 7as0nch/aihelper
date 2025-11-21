import { Clock, Calendar, User } from 'lucide-vue-next';

export interface MentionOption {
    id: string;
    label: string;
    value: string;
    description?: string;
}

export interface MentionType {
    key: string;
    label: string;
    icon: any;
    trigger: string; // The text after @ that triggers this type, e.g. "time" for @time
    hasSubMenu: boolean;
    fetchOptions?: (query?: string) => Promise<MentionOption[]>;
}

// Mock API for users
const mockUsers: MentionOption[] = [
    { id: '1', label: '张三', value: 'user:1', description: '前端开发' },
    { id: '2', label: '李四', value: 'user:2', description: '产品经理' },
    { id: '3', label: '王五', value: 'user:3', description: 'UI设计师' },
    { id: '4', label: '赵六', value: 'user:4', description: '测试工程师' },
];

export const mentionTypes: MentionType[] = [
    {
        key: 'time',
        label: '快速时间',
        icon: Clock,
        trigger: 'time',
        hasSubMenu: true,
        fetchOptions: async () => [
            { id: 'today', label: '今天', value: 'time:today' },
            { id: 'tomorrow', label: '明天', value: 'time:tomorrow' },
            { id: 'next_week', label: '下周', value: 'time:next_week' },
        ]
    },
    {
        key: 'range',
        label: '时间范围',
        icon: Calendar,
        trigger: 'range',
        hasSubMenu: true,
        fetchOptions: async () => [
            { id: 'last_7_days', label: '过去7天', value: 'range:last_7_days' },
            { id: 'last_30_days', label: '过去30天', value: 'range:last_30_days' },
            { id: 'this_month', label: '本月', value: 'range:this_month' },
        ]
    },
    {
        key: 'user',
        label: '用户',
        icon: User,
        trigger: 'user',
        hasSubMenu: true,
        fetchOptions: async (query) => {
            // Simulate API call
            await new Promise(resolve => setTimeout(resolve, 300));
            if (!query) return mockUsers;
            return mockUsers.filter(u => u.label.includes(query));
        }
    }
];
