import request from '@/utils/request';

export function sendMessage(data: any) {
    return request({
        url: '/chat/send',
        method: 'post',
        data,
    });
}
