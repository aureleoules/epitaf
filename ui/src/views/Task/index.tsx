import React, { useState } from 'react';
import { Task } from '../../types/task';
import styles from './task.module.scss';
import {ReactComponent as CalendarIcon} from '../../assets/svg/calendar.svg';
import {ReactComponent as ClockIcon} from '../../assets/svg/clock.svg';
import Button from '../../components/Button';
import ReactMarkdown from 'react-markdown'
import dayjs from 'dayjs';
import Input from '../../components/Input';
import Select from '../../components/Select';
import { DatePicker, MuiPickersUtilsProvider } from '@material-ui/pickers';
import DateDayjsUtils from '@date-io/dayjs';

import Client from '../../services/client';
import { useTranslation } from 'react-i18next';
import Checkbox from '../../components/Checkbox';
import {ReactComponent as LinkIcon} from '../../assets/svg/link.svg';
import { copy, capitalize, getUser, getSubjects } from '../../utils';

type Props = {
    task: Task
    new?: boolean
    close?: any
}
export default function(props: Props) {
    const {t} = useTranslation();

    const task = props.task;
    const [edit, setEdit] = useState<boolean>(false);
    const [content, setContent] = useState<string>("");
    const [subject, setSubject] = useState<string>("algorithmics");
    const [title, setTitle] = useState<string>("");
    const [due_date, setDueDate] = useState<Date>(dayjs(new Date()).add(24, 'hour').toDate());
    const [global, setGlobal] = useState<boolean>(false);

    const [promotion, setPromotion] = useState<number>(new Date().getFullYear()+5);
    const [classroom, setClass] = useState<string>("");
    const [region, setRegion] = useState<string>("Paris");
    const [semester, setSemester] = useState<string>("");

    function startEdit() {
        setContent(task.content!);
        setSubject(task.subject!);
        setTitle(task.title!);
        setDueDate(task.due_date!);
        setGlobal(task.global!);
        setPromotion(task.promotion!);
        setClass(task.class!);
        setRegion(task.region!);
        setSemester(task.semester!);
        setEdit(true);
    }

    function save() {
        const ta: Task = {
            content,
            subject,
            title,
            due_date,
            global,
            region,
            promotion,
            semester,
            class: classroom
        };

        if(!props.new) ta.short_id = task.short_id;

        if(props.new) {
            Client.Tasks.create(ta).then(id => {
                if(props.close) props.close();
            }).catch(err => {
                if(err) throw err;
            });
        } else {
            Client.Tasks.save(ta).then(id => {
                if(props.close) props.close();
            }).catch(err => {
                if(err) throw err;
            });
        }
    }

    function copyURL() {
        copy("https://" + window.location.host + "/tasks/" + task.short_id);   
    }

    function deleteTask() {
        if(!window.confirm(t('Delete task?'))) return;
    
        Client.Tasks.delete(task.short_id!).then(() => {
            // TODO fetch new list
            // avoid cache
            window.location.reload();
        }).catch(err => {
            if(err) throw err;
        });
    }
    
    return (
        <div className={styles.task}>
            {(!edit && !props.new) && <><div className={styles.header}>
                <h1>{task.title}</h1>
                <Button onClick={startEdit} title={t('Edit')}/>
                <LinkIcon className={styles.copy} onClick={copyURL}/>
            </div>
            <div className={styles.subheader}>
                <p className={styles.subject}>
                    {t(capitalize(task.subject!))}
                </p>
                <p className={styles.author}>
                    {t('Updated by')} <span>{task.updated_by}</span> {dayjs(task.updated_at).fromNow()}
                </p>
            </div>
            <div className={styles["content-container"]}>
                <ReactMarkdown allowedTypes={["break", "link", "text", "code", "blockquote", "strong", "emphasis", "list", "listItem", "paragraph", "thematicBreak", "heading"]} className={styles.content}>
                    {task.content!}
                </ReactMarkdown>
            </div>

            <div className={styles.columns}>
                <div className={styles.column}>
                    <h3>{t('Due date')}</h3>
                    <p>
                        <CalendarIcon/> {dayjs(task.due_date).format("DD MMMM")}
                    </p>
                </div>
                <div className={styles.column}>
                    <h3>{t('Updated')}</h3>
                    <p>
                        <ClockIcon/> {dayjs(task.updated_at).fromNow()}
                    </p>
                </div>
                <div className={styles.column}>
                    <h3>{t('Created')}</h3>
                    <p>
                        <ClockIcon/> {dayjs(task.created_at).fromNow()}
                    </p>
                </div>
            </div>
            </>}

            {(edit || props.new) && <div className={styles.edit}>
                <h1>{t('Add a task')}</h1>
                <Select value={subject} onChange={(e:any) => setSubject(e.target.value)} title={t("Subject")}>
                    {getSubjects(getUser().teacher)
                        .sort((a, b) => t(a.display_name).localeCompare(t(b.display_name)))
                        .map((s, i) => <option key={i} value={s.name}>
                        {t(s.display_name)}
                    </option>)}
                </Select>
                <Input value={title} onChange={(e: any) => setTitle(e.target.value)} placeholder={t('Title')}/>
                <Input
                    onChange={(e:any) => setContent(e.target.value)}
                    multiline 
                    value={content} 
                    placeholder={t("Content")} 
                    rows={10}
                />
                <MuiPickersUtilsProvider utils={DateDayjsUtils} locale={"fr"}>
                    <DatePicker
                        variant="dialog"
                        label={t('Due date')}
                        value={due_date}
                        autoOk
                        onChange={d => setDueDate(d?.toDate()!)}
                    />
                </MuiPickersUtilsProvider>
                <Checkbox disabled={!props.new} checked={global} onChange={(e:any) => setGlobal(e.target.checked)} title={t('Promotion')}/>

                {getUser().teacher && <>
                    <div className={styles.classinfos}>
                        <Input value={promotion} onChange={(e: any) => setPromotion(parseInt(e.target.value))} type="number" placeholder={t('Promotion')}/>
                        <Input value={semester} onChange={(e: any) => setSemester(e.target.value)} disabled={global} placeholder={t('Semester')}/>
                        <Input value={classroom} onChange={(e: any) => setClass(e.target.value)} disabled={global} placeholder={t('Class')}/>
                    </div>
                    <Select value={region} onChange={(e: any) => setRegion(e.target.value)} disabled={global} title={t('Region')}>
                        {["Paris", "Lyon", "Strasbourg", "Rennes", "Toulouse"].map((r, i) => <option value={r}>{r}</option>)}
                    </Select>
                </>}
                <div className={styles.actions + " " + (props.new ? styles.new : "")}>
                    <Button className={styles.save} disabled={!subject || !title || !content} onClick={save} title={props.new ? t("Create"): t("Save")}/>
                    {!props.new && <Button color="red" onClick={deleteTask} title={t('Delete')}/>}
                </div>
            </div>}
        </div>
    )
}