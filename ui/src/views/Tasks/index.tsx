import React, { useEffect, useState } from 'react';
import dayjs from 'dayjs';
import styles from './tasks.module.scss';
import relativeTime from 'dayjs/plugin/relativeTime';
import { Task } from '../../types/task';
import Button from '../../components/Button';
import {ReactComponent as PlusIcon} from '../../assets/svg/plus.svg';
import {ReactComponent as UsersIcon} from '../../assets/svg/users.svg';
import {ReactComponent as UserIcon} from '../../assets/svg/user.svg';
import {ReactComponent as CheckIcon} from '../../assets/svg/check.svg';
import { Link } from 'react-router-dom';

import TaskView from '../Task';
import Modal from '../../components/Modal';
import Client from '../../services/client';
import { useTranslation } from 'react-i18next';
import { capitalize, deleteFilters, getSubjects, getUser, loadFilters, saveFilters } from '../../utils';
import { RotateSpinner  } from "react-spinners-kit";
import history from '../../history';
import Datetime from 'react-datetime';
import "react-datetime/css/react-datetime.css";
import { Filters } from '../../types/filters';
import Select from '../../components/Select';
import Checkbox from '../../components/Checkbox';

dayjs.extend(relativeTime)

const subjects = getSubjects(getUser().teacher);

type Props = {
    location: any
}
export default function(props: any) {
    const {t} = useTranslation();

    
    const [tasks, setTasks] = useState<Array<Task>>(new Array<Task>());
    const [open, setOpen] = useState<boolean>(false);
    const [task, setTask] = useState<Task | null>(null);
    const [created_task, setCreateTask] = useState<boolean>(false);
    const [fetched, setFetched]= useState<boolean>(false);
    

    // Filters
    const [filters, setFilters] = useState<boolean>(false);

    // Defaults
    const f = loadFilters();

    const [startDate, setStartDate] = useState<Date>(new Date());
    const [endDate, setEndDate] = useState<Date>(dayjs(new Date()).add(2, 'month').toDate());
    const [subject, setSubject] = useState<string | undefined>(f.subject);
    const [visibility, setVisibility] = useState<string | undefined>(f.visibility);
    const [status, setStatus] = useState<string | undefined>(f.status);

    const [activeFilters, setActiveFilters] = useState<boolean>(f.active! || false);
    
    function getFilters(): Filters {
        let completed;
        if(status === "completed") completed = true;
        else if(status === "todo") completed = false;
        else completed = undefined;

        return {
            start_date: startDate,
            end_date: endDate,
            status: status || undefined,
            completed: completed,
            visibility: visibility || undefined,
            subject: subject || undefined,
            active: activeFilters || false 
        };
    }
    
    function fetchTasks() {
        const filters = activeFilters ? getFilters() : undefined;

        Client.Tasks.list(filters).then(tasks => {
            setTasks(tasks);
            setFetched(true);
        });
    }

    useEffect(() => {
        fetchTasks();

        if(!props.match) return;
        const id = props.match.params.id;
        if(id) {
            Client.Tasks.get(id).then(t => {
                openTask(t);
            }).catch(err => {
                if(err) throw err;
            });
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [props.match, startDate, endDate, status, subject, visibility, activeFilters]);

    useEffect(() => {
        const filters = getFilters();
        saveFilters(filters);
        
     // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [startDate, endDate, status, subject, visibility, activeFilters]);

    function resetFilters() {
        setStartDate(new Date());
        setEndDate(dayjs(new Date()).add(2, 'month').toDate());
        setSubject("");
        setVisibility("");
        setStatus("");
    }

    function openTask(t: Task) {
        setTask(t);
        setOpen(true);
    }
    
    function closeModal() {
        setOpen(false);
        let pathname = history.location.pathname.split("/")[1];
        if(pathname === 't') pathname = 'tasks';
        history.push(`/${pathname}`);
        setCreateTask(false);
    }

    function createTask() {
        setOpen(true);
        setCreateTask(true);
    }
    
    function getIcon(subject: string): string {
        for(let i = 0; i < subjects.length; i++) {
            if(subjects[i].name === subject) return subjects[i].icon;
        }
        return "document.svg";
    }

    function getSubject(subject: string): string {
        for(let i = 0; i < subjects.length; i++) {
            if(subjects[i].name === subject) return subjects[i].display_name;
        }
        return subject;
    }

    function sortTasks(tasks: Array<Task>): Array<Array<Task>> {
        const sorted = new Array<Array<Task>>();
        
        for(let i = 0; i < tasks.length; i++) {
            const task = tasks[i];
            const date = new Date(task.due_date!);
            let found = false;
            for(let j = 0; j < sorted.length; j++) {
                const cmp = new Date(sorted[j][0].due_date!);
                if(cmp.getDate() === date.getDate() && cmp.getMonth() === date.getMonth() && cmp.getFullYear() === date.getFullYear()) {
                    sorted[j].push(task);
                    found = true;
                    break;
                }
            }
            if(!found) sorted.push([task]);
        }

        sorted.sort((a, b) => dayjs(b[0].due_date!).isBefore(dayjs(a[0].due_date!)) ? 1 : -1);
        return sorted;
    }
    
    
    function complete(e: any, ta: Task) {
        e.stopPropagation();
        e.preventDefault();
        if(ta.completed) {
            Client.Tasks.incomplete(ta.short_id!).then(() => fetchTasks()).catch(err => {
                if(err) throw err;
            });
        } else {
            Client.Tasks.complete(ta.short_id!).then(() => fetchTasks()).catch(err => {
                if(err) throw err;
            });
        }
    }

    return (
        <div className={styles.tasks}>
            <div className={styles.header}>
                <h1>{t('Tasks')}</h1>
                <Button onClick={() => setFilters(!filters)} className={styles.filtersbtn} title={t('Filters')}>
                    {activeFilters && <span className={styles.bullet}>1</span>}
                </Button>
                <Button className={styles.addbtn} onClick={createTask} icon={PlusIcon} title={t('New')}/>
            </div>
            
            {filters && <div className={styles.filters}>
                <div className={styles.row}>
                    <div className={styles.dateinput}>
                        <p>{t('Start date')}</p>
                        <Datetime 
                            inputProps={{
                                placeholder: t('Start date')
                            }} 
                            dateFormat="DD MMMM YYYY" 
                            timeFormat={false} 
                            className={styles.datepicker} 
                            value={startDate}
                            onChange={(value: any) => setStartDate(value.toDate())}
                            closeOnSelect
                        />
                    </div>

                    <div className={styles.dateinput}>
                        <p>{t('End date')}</p>
                        <Datetime 
                            inputProps={{
                                placeholder: t('End date')
                            }} 
                            dateFormat="DD MMMM YYYY" 
                            timeFormat={false} 
                            className={styles.datepicker} 
                            initialValue={dayjs(new Date()).add(1, 'month').toDate()}
                            value={endDate}
                            onChange={(value: any) => setEndDate(value.toDate())}
                            closeOnSelect
                        />
                    </div>
                </div>
                <div className={styles.row}>
                    <Select 
                        className={styles.visibility} 
                        value={visibility} 
                        onChange={(e: any) => setVisibility(e.target.value)}    
                        title={t('Visibility')}>
                        <option value={""}>{t('Any')}</option>
                        <option value={'self'}>{t('Me')}</option>
                        <option value={'students'}>{t('Students')}</option>
                        <option value={'class'}>{t('Classe') + (!getUser().teacher ? ` (${getUser().class})` : "")}</option>
                        <option value={'promotion'}>{t('Promotion') + (!getUser().teacher ? ` (${getUser().promotion})` : "")}</option>
                    </Select>
                    <Select value={subject} onChange={(e:any) => setSubject(e.target.value)} title={t("Subject")}>
                        <option value={""}>{t('Any')}</option>
                        {getSubjects(getUser().teacher)
                            .sort((a, b) => t(a.display_name).localeCompare(t(b.display_name)))
                            .map((s, i) => <option key={i} value={s.name}>
                            {t(s.display_name)}
                        </option>)}
                    </Select>
                    <Select value={status} onChange={(e:any) => setStatus(e.target.value)} title={t("Status")}>
                        <option value={""}>{t('Any')}</option>
                        <option value={"completed"}>{t('Completed')}</option>
                        <option value={"todo"}>{t('To do')}</option>
                    </Select> 
                </div>
                <div className={styles.row}>
                    {/* <Button title={t('Apply')}/> */}
                    <Checkbox 
                        color="green" 
                        className={styles.checkbox} 
                        title={t('Active')}
                        checked={activeFilters} 
                        onChange={(e: any) => setActiveFilters(e.target.checked)}
                    />
                    <Button onClick={resetFilters} className={styles.reset} color="red" title={t('Reset')}/>
                </div>
            </div>}


            {!fetched && <div style={{position: "absolute", left: "40%", top: "45%"}}>
                <RotateSpinner size={50} color="var(--primary)"/>
            </div>}

            {(fetched && tasks.length === 0) && <h4>{t('Nothing to do.')}</h4>}
            {sortTasks(tasks).map((l, i) => {
                return (
                    <div key={i}>
                        <h2>{capitalize(dayjs(l[0].due_date).format("dddd DD MMMM"))}</h2>
                        {l.sort((a, b) => dayjs(a.updated_at!).isBefore(dayjs(b.updated_at!)) ? 1 : -1).map((ta, i) => (
                            <Link to={"#"} onClick={() => openTask(ta)} key={i} className={[styles.task, ta.completed ? styles.completed: ""].join(" ")}>
                                <div className={styles["icon-container"]}>
                                    <img alt={ta.subject} width={30} src={require('../../assets/svg/subjects/' + getIcon(ta.subject!))}/>
                                </div>
                                <div className={styles.content}>
                                    <h4>{t(capitalize(getSubject(ta.subject!)))}</h4>
                                    <p>{ta.title?.substr(0, 50)} · {t('Edited')} {dayjs(ta.updated_at).fromNow()}</p>
                                </div>
                                <div className={styles.infos}>
                                    <span onClick={(e: any) => complete(e, ta)} className={[styles.complete, ta.completed ? styles.completed : styles.incompleted].join(" ")}><CheckIcon/></span>
                                    {ta.visibility === "promotion" && <span className={styles.group + " " + styles.promotion}>
                                        {ta.promotion}
                                    </span>}
                                    {ta.visibility === "class" && <span className={styles.group + " " + styles.class}>
                                        {ta.class || ta.region}
                                    </span>}
                                    {ta.visibility === "self" && <span className={styles.group + " " + styles.icon + " " + styles.self}>
                                        <UserIcon/>
                                    </span>}
                                    {ta.visibility === "students" && <span className={styles.group + " " + styles.icon + " " + styles.students}>
                                        <UsersIcon/>
                                    </span>}
                                </div>
                            </Link>
                        ))}
                    </div>
                );
            })}
            {/* TODO: show more */}
            {/* <div>
                <Button className={styles.showmore} title={t('Show more')}/>
            </div> */}
            {open && <Modal close={closeModal}>
                <TaskView close={() => {
                    fetchTasks();
                    closeModal();
                }} new={created_task} task={task!}/>
            </Modal>}
        </div>
    )
}