import React from 'react';
import styles from './select.module.scss';

type Props = {
    title?: string
    children: any
    value?: any
    onChange?: any
    disabled?: boolean
    className?: string
}
export default function(props: Props) {
    return (
        <div className={[styles.select, props.className].join(" ")}>
            <p>{props.title}</p>
            <select disabled={props.disabled} value={props.value} onChange={props.onChange}>
                {props.children}
            </select>
        </div>
    )
}