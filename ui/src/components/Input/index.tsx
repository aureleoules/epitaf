import React from 'react';

import styles from './input.module.scss';

type Props = {
    placeholder?: string
    value?: any
    onChange?: any
    multiline?: boolean
    rows?: number
    title?: string
    className?: string
    type?: string
    disabled?: boolean
}
export default function(props: Props) {

    const p = {
        placeholder: props.placeholder, 
        value: props.value,
        onChange: props.onChange, 
        className: props.className,
        disabled: props.disabled
    }
    
    let el = props.multiline ? <textarea rows={props.rows} {...p}/> : <input type={props.type} {...p}/>;
    return <div className={styles.input}>
        <p className={styles.p}>{props.title || props.placeholder}</p>
        {el}
    </div>
}