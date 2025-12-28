import React, { useEffect, useState } from 'react';
// Cache Buster: 2025-12-28-v2
import api from '../services/api';
import useAuth from '../hooks/useAuth';
import { Tag, Home, Plus, Trash2 } from 'lucide-react';

const Sidebar = ({ selectedTagId, onSelectTag }) => {
    const { user } = useAuth();
    const [tags, setTags] = useState([]);
    const [isCreatingTag, setIsCreatingTag] = useState(false);
    const [newTagName, setNewTagName] = useState('');

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

    const handleCreateTag = async (e) => {
        e.preventDefault();
        if (!newTagName.trim()) return;

        try {
            await api.post('/tags', {
                name: newTagName
            });
            setNewTagName('');
            setIsCreatingTag(false);
            fetchTags(); // Refresh list
        } catch (err) {
            console.error("Failed to create tag", err);
            alert("Failed to create tag");
        }
    };

    const handleDeleteTag = async (e, tagId) => {
        e.stopPropagation(); // Don't select the tag
        if (!window.confirm("Delete this tag?")) return;

        try {
            await api.delete(`/tags/${tagId}`);
            if (selectedTagId === tagId) onSelectTag(null); // Deselect if deleted
            fetchTags();
        } catch (err) {
            console.error("Failed to delete tag", err);
            alert("Failed to delete tag");
        }
    };

    return (
        <aside style={{ width: '250px', background: 'white', borderRight: '1px solid var(--border-color)', height: 'calc(100vh - 80px)', overflowY: 'auto', padding: '1rem', display: 'flex', flexDirection: 'column' }}>
            <nav style={{ flex: 1 }}>
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

                <div style={{ marginTop: '2rem', marginBottom: '1rem', paddingLeft: '0.75rem', fontSize: '0.875rem', color: 'var(--text-muted)', fontWeight: 'bold', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <span>TAGS</span>
                    <button
                        onClick={() => setIsCreatingTag(!isCreatingTag)}
                        style={{ padding: '4px', borderRadius: '4px', color: 'var(--primary-color)' }}
                        title="Create Tag"
                    >
                        <Plus size={16} />
                    </button>
                </div>

                {isCreatingTag && (
                    <form onSubmit={handleCreateTag} style={{ marginBottom: '1rem', padding: '0 0.75rem' }}>
                        <input
                            autoFocus
                            type="text"
                            placeholder="Tag name..."
                            value={newTagName}
                            onChange={(e) => setNewTagName(e.target.value)}
                            style={{ width: '100%', padding: '0.5rem', borderRadius: '4px', border: '1px solid var(--primary-color)', fontSize: '0.875rem' }}
                            onBlur={() => !newTagName && setIsCreatingTag(false)}
                        />
                    </form>
                )}

                {tags.map(tag => (
                    <div
                        key={tag.tag_id}
                        onClick={() => onSelectTag(tag.tag_id)}
                        className="tag-item" // For potential CSS hover effects if needed
                        style={{
                            display: 'flex', alignItems: 'center', gap: '0.75rem', padding: '0.75rem',
                            borderRadius: 'var(--radius-md)',
                            cursor: 'pointer',
                            background: selectedTagId === tag.tag_id ? 'var(--bg-secondary)' : 'transparent',
                            color: selectedTagId === tag.tag_id ? 'var(--primary-color)' : 'inherit',
                            position: 'relative',
                            group: 'true'
                        }}
                        onMouseEnter={e => e.currentTarget.querySelector('.delete-btn').style.opacity = 1}
                        onMouseLeave={e => e.currentTarget.querySelector('.delete-btn').style.opacity = 0}
                    >
                        <Tag size={18} />
                        <span style={{ flex: 1, whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>{tag.name}</span>

                        <button
                            className="delete-btn"
                            onClick={(e) => handleDeleteTag(e, tag.tag_id)}
                            style={{
                                opacity: 0, transition: 'opacity 0.2s',
                                border: 'none', background: 'transparent', color: '#ef4444',
                                padding: '4px', cursor: 'pointer'
                            }}
                            title="Delete Tag"
                        >
                            <Trash2 size={14} />
                        </button>
                    </div>
                ))}
            </nav>
        </aside>
    );
};

export default Sidebar;
