import React, { useEffect, useState } from 'react';
import dayjs from 'dayjs';
import styles from './tasks.module.scss';
import relativeTime from 'dayjs/plugin/relativeTime';
import { Task } from '../../types/task';
import Button from '../../components/Button';
import {ReactComponent as PlusIcon} from '../../assets/svg/plus.svg';
import {ReactComponent as UsersIcon} from '../../assets/svg/users.svg';
import {ReactComponent as UserIcon} from '../../assets/svg/user.svg';
import { Link } from 'react-router-dom';

import TaskView from '../Task';
import Modal from '../../components/Modal';
import Client from '../../services/client';
import { useTranslation } from 'react-i18next';
import { capitalize, getSubjects, getUser } from '../../utils';
import { RotateSpinner  } from "react-spinners-kit";
import history from '../../history';

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
    
    function fetchTasks() {
        Client.Tasks.list().then(tasks => {
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
    }, [props.match]);

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
    
    
    return (
        <div className={styles.tasks}>
            <div className={styles.header}>
                <h1>{t('Tasks')}</h1>
                <Button onClick={createTask} icon={PlusIcon} title={t('New')}/>
            </div>

            {!fetched && <div style={{position: "absolute", left: "40%", top: "45%"}}>
                <RotateSpinner size={50} color="var(--primary)"/>
            </div>}

            {(fetched && tasks.length === 0) && <h4>{t('Nothing to do.')}</h4>}
            {sortTasks(tasks).map((l, i) => {
                return (
                    <div key={i}>
                        <h2>{capitalize(dayjs(l[0].due_date).format("dddd DD MMMM"))}</h2>
                        {l.reverse().map((ta, i) => (
                            <Link to={"#"} onClick={() => openTask(ta)} key={i} className={styles.task}>
                                <div className={styles["icon-container"]}>
                                    <img alt={ta.subject} width={30} src={require('../../assets/svg/subjects/' + getIcon(ta.subject!))}/>
                                </div>
                                <div className={styles.content}>
                                    <h4>{t(capitalize(getSubject(ta.subject!)))}</h4>
                                    <p>{ta.title?.substr(0, 50)} · {t('Edited')} {dayjs(ta.updated_at).fromNow()}</p>
                                </div>
                                <div className={styles.infos}>
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
            {open && <Modal close={closeModal}>
                <TaskView close={() => {
                    fetchTasks();
                    closeModal();
                }} new={created_task} task={task!}/>
            </Modal>}
        </div>
    )
}