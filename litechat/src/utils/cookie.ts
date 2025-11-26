export const setToken = (token: string) => {
    // 7 days in seconds
    const maxAge = 7 * 24 * 60 * 60;
    document.cookie = `token=${token}; path=/; max-age=${maxAge}; SameSite=Strict`;
};

export const getToken = (): string | null => {
    const name = 'token=';
    const decodedCookie = decodeURIComponent(document.cookie);
    const ca = decodedCookie.split(';');
    for (let i = 0; i < ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) === ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) === 0) {
            return c.substring(name.length, c.length);
        }
    }
    return null;
};

export const removeToken = () => {
    document.cookie = 'token=; path=/; max-age=0; SameSite=Strict';
};
