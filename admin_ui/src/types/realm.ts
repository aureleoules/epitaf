import { base } from './base';

export type Realm = base & {
    uuid?: string
    name?: string
    slug?: string
    url?: string
    website_url?: string
}