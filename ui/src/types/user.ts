import { base } from './base';

export type User = base & {
    name?: string
    login?: string
    promotion?: number
    class?: string
    region?: string
    semester?: string
    email?: string
    teacher?: boolean
}