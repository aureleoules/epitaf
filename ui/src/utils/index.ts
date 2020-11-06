import { User } from "../types/user";
import subjects from '../assets/data/subjects.json';
import { Subject } from "../types/subject";

export const getQueryVariable = (variable: string) => {
    const query = window.location.search.substring(1);
    const vars = query.split('&');
    for (let i = 0; i < vars.length; i++) {
        const pair = vars[i].split('=');
        if (decodeURIComponent(pair[0]) === variable) {
            return decodeURIComponent(pair[1]);
        }
    }
    return null;
}

export const capitalize = (str: string): string => {
    return str.charAt(0).toUpperCase() + str.slice(1);
}

export const isLoggedIn = (): boolean => {
    return !!localStorage.getItem("jwt");
}

export const copy = (str: string) => {
    const input = document.createElement('input');
    input.setAttribute('value', str);
    input.style.position = "absolute";
    input.style.left = "-99999px";
    input.style.opacity = "0";
    document.body.appendChild(input);
    input.select();
    document.execCommand('copy');
    document.body.removeChild(input);
}

export const parseJwt = (token: string): any => {
    if(!token) return {};
    var base64Url = token.split('.')[1];
    var base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    var jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
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

export const getClassType = (semester: string): string => {
    if(!semester) return "";
    
    const n = parseInt(semester[1]);
    if(n <= 2) return "SUP";
    if(n <= 4) return "SPE";
    return "ING";
}

export const getSubjects = (all?: boolean): Array<Subject> => {
    if(all) return subjects;
    
    const user = getUser();
    const type = getClassType(user.semester!);

    return subjects.filter(s => s.classes.includes(type));
}