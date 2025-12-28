import React, { useEffect, useState } from 'react';
import api from '../services/api';
import { Tag, Home } from 'lucide-react';

const Sidebar = ({ selectedTagId, onSelectTag }) => {
    const [tags, setTags] = useState([]);

    useEffect(() => {
        fetchTags();
    }, []);

    const fetchTags = async () => {
        try {
            const res = await api.get('/tags?tag_id=0&page_size=100');
            setTags(res.data || []);
        } catch (err) {
            console.error("Failed to fetch tags", err);
        }
    };

    return (
        <aside style={{ width: '250px', background: 'white', borderRight: '1px solid var(--border-color)', height: 'calc(100vh - 80px)', overflowY: 'auto', padding: '1rem' }}>
            <nav>
                <div
                    onClick={() => onSelectTag(null)}
                    style={{
                        display: 'flex', alignItems: 'center', gap: '0.75rem', padding: '0.75rem',
                        borderRadius: 'var(--radius-md)',
                        cursor: 'pointer',
                        background: selectedTagId === null ? 'var(--bg-secondary)' : 'transparent',
                        color: selectedTagId === null ? 'var(--primary-color)' : 'inherit',
                        fontWeight: selectedTagId === null ? 'bold' : 'normal'
                    }}
                >
                    <Home size={20} />
                    <span>All Notes</span>
                </div>

                <div style={{ marginTop: '2rem', marginBottom: '1rem', paddingLeft: '0.75rem', fontSize: '0.875rem', color: 'var(--text-muted)', fontWeight: 'bold' }}>
                    TAGS
                </div>

                {tags.map(tag => (
                    <div
                        key={tag.tag_id}
                        onClick={() => onSelectTag(tag.tag_id)}
                        style={{
                            display: 'flex', alignItems: 'center', gap: '0.75rem', padding: '0.75rem',
                            borderRadius: 'var(--radius-md)',
                            cursor: 'pointer',
                            background: selectedTagId === tag.tag_id ? 'var(--bg-secondary)' : 'transparent',
                            color: selectedTagId === tag.tag_id ? 'var(--primary-color)' : 'inherit'
                        }}
                    >
                        <Tag size={18} />
                        <span>{tag.name}</span>
                    </div>
                ))}
            </nav>
        </aside>
    );
};

export default Sidebar;
