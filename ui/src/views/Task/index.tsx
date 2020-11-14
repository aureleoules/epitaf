import React, { useEffect, useRef, useState } from 'react';
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
import {ReactComponent as LinkIcon} from '../../assets/svg/link.svg';
import { copy, capitalize, getUser, getSubjects } from '../../utils';
import { IDictionary } from '../../types/dictionnary';
import TagsInput from 'react-tagsinput';
import 'react-tagsinput/react-tagsinput.css';
import { User } from '../../types/user';

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

    const [promotion, setPromotion] = useState<number>(new Date().getFullYear()+5);
    const [classroom, setClass] = useState<string>("");
    const [region, setRegion] = useState<string>("");
    const [semester, setSemester] = useState<string>("");

    const [members, setMembers] = useState<Array<string>>(new Array<string>());
    const [classes, setClasses] = useState<IDictionary<any>>();

    const [searchedUsers, setSearchedUsers] = useState<Array<User>>(new Array<User>());

    const [visibility, setVisibility] = useState<string>("self");

    function startEdit() {
        setContent(task.content!);
        setSubject(task.subject!);
        setTitle(task.title!);
        setDueDate(task.due_date!);
        setVisibility(task.visibility!);
        setPromotion(task.promotion!);
        setClass(task.class!);
        setRegion(task.region!);
        setSemester(task.semester!);
        setEdit(true);
        setMembers(task.members || []);
    }

    function save() {
        const ta: Task = {
            content,
            subject,
            members,
            title,
            due_date,
            visibility,
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
    
    function fetchClasses() {
        Client.Classes.list().then(classes => {
            setClasses(classes);

            setTimeout(() => {
                const p = parseInt(Object.keys(classes!)[0]);
                updatePromotion(p, classes);
            }, 10);

        }).catch(err => {
            if(err) throw err;
        });
    }

    useEffect(() => {
        fetchClasses();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    function updatePromotion(p: number, c: IDictionary<any> = classes!) {
        setPromotion(p);
        const s = Object.keys(c![p])[0];
        setSemester(s);
        const r = Object.keys(c![p][s])[0];
        setRegion(r);

        setClass(c![p][s][r][0]);
    }

    function updateSemester(s: string, c: IDictionary<any> = classes!) {
        setSemester(s);
        const r = Object.keys(c![promotion][s])[0];
        setRegion(r);

        setClass(c![promotion][s][r][0]);
    }

    function updateRegion(r: string, c: IDictionary<any> = classes!) {
        setRegion(r);
        setClass(c![promotion][semester][r][0]);
    }

    function searchUser(q: string) {
        if(q.length === 0) setSearchedUsers(new Array<User>());
        if(q.length < 2) return;
        
        Client.Users.search(q).then(users => {
            setSearchedUsers(users);
        }).catch(err => {
            if(err) throw err;
        });;
    }

    const tagRef = useRef(null);
    function onKeyEnter(e: any) {
        if(e.key === "Enter") {
            const filtered = searchedUsers.filter((el) => !members.includes(el.login!));
            if(!filtered || filtered.length < 1) return;
            addMember(filtered[0].login!);
        }
    }

    function addMember(login: string) {
        if(members.includes(login)) return;
        setMembers(m => [...m, login]);
        setTimeout(() => {
            const u = (tagRef.current! as any);
            u.clearInput()
            u.focus();
            setSearchedUsers([]);
        }, 1);
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
                
                {(props.new || getUser().login === task.created_by_login) && <>
                    <Select value={visibility} onChange={(e: any) => setVisibility(e.target.value)} title={t('Visibility')}>
                        <option value={'self'}>{t('Me')}</option>
                        <option value={'students'}>{t('Students')}</option>
                        <option value={'class'}>{t('Classe') + (!getUser().teacher ? ` (${getUser().class})` : "")}</option>
                        {!getUser().teacher && <option value={'promotion'}>{t('Promotion') + ` (${getUser().promotion})`}</option>}
                    </Select>
                    {getUser().teacher && classes && <>
                        <div className={styles.classinfos}>
                            <Select value={promotion} onChange={(e: any) => {
                                updatePromotion(e.target.value);
                            }} title={t('Promotion')}>
                                {Object.keys(classes!).map((r: string, i: number) => <option value={r}>{r}</option>)}
                            </Select>
                            <Select disabled={!promotion} value={semester} onChange={(e: any) => {
                                updateSemester(e.target.value);
                            }} title={t('Semester')}>
                                {promotion && Object.keys(classes[promotion]!).map((s: string, i: number) => <option value={s}>{s}</option>)}
                            </Select>
                            <Select disabled={!semester || visibility === "promotion"} value={region} onChange={(e: any) => {
                                updateRegion(e.target.value);
                            }} title={t('Region')}>
                                {((semester && visibility !== "promotion") && classes[promotion][semester]) ? (Object.keys(classes[promotion][semester]).map((r: string, i: number) => <option value={r}>{r}</option>)) : null}
                            </Select>
                        </div>
                        {<Select value={classroom} onChange={(e: any) => setClass(e.target.value)} disabled={visibility === "promotion" || !region} title={t('Class')}>
                            {(region && visibility !== "promotion") && classes[promotion][semester][region].map((r: string, i: number) => <option value={r}>{r === "" ? t('All') : r}</option>)}
                        </Select>}
                    </>}

                    {visibility === 'students' &&
                        <>
                            <TagsInput 
                                className={"react-tagsinput"}
                                inputProps={{
                                    placeholder: t('Students')
                                }} 
                                value={members}
                                ref={tagRef}
                                renderInput={props => <input {...props} value={props.value} onKeyDown={e => {
                                    if(e.key === "Enter") {
                                        if(searchedUsers.length < 1) return;
                                        onKeyEnter(e);
                                    } else props.onKeyDown(e);
                                }} onChange={e => {
                                    searchUser(e.target.value);
                                    props.onChange(e);
                                }} />}
                                onChange={members => setMembers(members)}
                            />
                            {searchedUsers.length > 0 && <div className={styles.users}>
                                {searchedUsers.filter((el) => !members.includes(el.login!)).map((u, i) => (
                                    <div onClick={() => {
                                        addMember(u.login!)
                                    }} className={styles.user}>
                                        <p>{u.name}</p>
                                    </div>
                                ))}
                            </div>}
                        </>}
                </>}
                

                <div className={styles.actions + " " + (props.new ? styles.new : "")}>
                    <Button className={styles.save} disabled={!subject || !title || !content} onClick={save} title={props.new ? t("Create"): t("Save")}/>
                    {!props.new && <Button color="red" onClick={deleteTask} title={t('Delete')}/>}
                </div>
            </div>}
        </div>
    )
}