import { base } from './base';

export type Group = base & {
    uuid?: string
    realm_id?: string
    usable?: boolean
    name?: string
    users?: number
    slug?: string
    parent_id?: string
    subgroups?: Array<Group>
    archived?: boolean
    archived_at?: Date
    active_at?: Date
}