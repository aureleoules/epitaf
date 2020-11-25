import React from 'react';
import styles from './checkbox.module.scss';

type Props = {
    title?: string
    onChange?: any
    checked?: boolean
    disabled?: boolean
    className?: string
    color?: "green" | "red"
}
export default function Checkbox(props: Props) {
    return (
        <div className={[styles.checkbox, props.disabled ? styles.disabled : "", props.className].join(" ")}>
            <label className={styles.container}>
                <p>
                    {props.title}
                </p>
                <input checked={props.checked} onChange={props.onChange} type="checkbox"/>
                <span className={[styles.checkmark, styles[props.color!]].join(" ")}></span>
            </label>
        </div>
    )
}