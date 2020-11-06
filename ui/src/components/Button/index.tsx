import React from 'react';
import styles from './button.module.scss';

type Props = {
    title: string
    onClick?: any
    icon?: any
    color?: "primary" | "green" | "red" | undefined
    large?: boolean
    disabled?: boolean
    className?: string
}
export default function(props: Props) {
    if(!props.color) props = Object.assign({}, props, {color: "primary"});
    
    return (
        <button 
            className={[
                styles.button, 
                styles[props.color!], 
                props.large ? styles.large : "",
                props.disabled ? styles.disabled : "",
                props.className
            ].join(" ")} 
            onClick={props.onClick}>
            {props.icon && <props.icon/>}
            {props.title}
        </button>
    )
}