import { User } from "../types/user";

export const isLoggedIn = (): boolean => {
    return !!localStorage.getItem("jwt");
}
export const parseJwt = (token: string): any => {
    if (!token) return {};
    var base64Url = token.split('.')[1];
    var base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    var jsonPayload = decodeURIComponent(atob(base64).split('').map(function (c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));

    return JSON.parse(jsonPayload);
};
export const logout = () => {
    localStorage.setItem("jwt", "");
    window.location.replace("/");
}

export const getUser = (): User => {
    return parseJwt(localStorage.getItem("jwt")!);
}