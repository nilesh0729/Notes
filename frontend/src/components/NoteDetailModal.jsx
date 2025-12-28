import React, { useState, useEffect } from 'react';
import api from '../services/api';
import { X, Tag as TagIcon, Plus, Trash2, Edit2, Save } from 'lucide-react';

const NoteDetailModal = ({ noteId, onClose, onNoteUpdated }) => {
    const [note, setNote] = useState(null);
    const [loading, setLoading] = useState(true);
    const [availableTags, setAvailableTags] = useState([]);
    const [isAddingTag, setIsAddingTag] = useState(false);

    // Edit Mode State
    const [isEditing, setIsEditing] = useState(false);
    const [editTitle, setEditTitle] = useState('');
    const [editContent, setEditContent] = useState('');

    useEffect(() => {
        fetchNoteDetails();
        fetchTags();
    }, [noteId]);

    const fetchNoteDetails = async () => {
        try {
            const res = await api.get(`/notes/${noteId}`);
            setNote(res.data);
            setEditTitle(res.data.title);
            setEditContent(res.data.content);
        } catch (err) {
            console.error("Failed to fetch note details", err);
        } finally {
            setLoading(false);
        }
    };

    const fetchTags = async () => {
        try {
            const res = await api.get('/tags?tag_id=0&page_size=100');
            setAvailableTags(res.data || []);
        } catch (err) {
            console.error("Failed to fetch tags", err);
        }
    };

    const handleAddTag = async (tagId) => {
        try {
            await api.post('/note_tags', {
                note_id: noteId,
                tag_id: tagId
            });
            setIsAddingTag(false);
            fetchNoteDetails(); // Refresh to show new tag
            onNoteUpdated && onNoteUpdated();
        } catch (err) {
            console.error("Failed to add tag", err);
            alert("Failed to add tag");
        }
    };

    const handleSave = async () => {
        try {
            await api.put(`/notes/${noteId}`, {
                title: editTitle,
                content: editContent
            });
            setIsEditing(false);
            fetchNoteDetails();
            onNoteUpdated && onNoteUpdated();
        } catch (err) {
            console.error("Failed to update note", err);
            alert("Failed to update note");
        }
    };

    const handleDelete = async () => {
        if (!window.confirm("Are you sure you want to delete this note?")) return;
        try {
            await api.delete(`/notes/${noteId}`);
            onNoteUpdated && onNoteUpdated();
            onClose();
        } catch (err) {
            console.error("Failed to delete note", err);
            alert("Failed to delete note");
        }
    };

    if (loading) return <div className="modal-overlay">Loading...</div>;
    if (!note) return null;

    const existingTagIds = note.tags ? note.tags.map(t => t.tag_id) : [];

    return (
        <div style={{
            position: 'fixed', top: 0, left: 0, right: 0, bottom: 0,
            background: 'rgba(0,0,0,0.5)', display: 'flex', justifyContent: 'center', alignItems: 'center',
            backdropFilter: 'blur(4px)', zIndex: 1000
        }}>
            <div style={{ background: 'white', padding: '2rem', borderRadius: '8px', width: '90%', maxWidth: '800px', maxHeight: '90vh', overflowY: 'auto', position: 'relative' }}>
                <button
                    onClick={onClose}
                    style={{ position: 'absolute', top: '1rem', right: '1rem', background: 'transparent', border: 'none', cursor: 'pointer', color: '#666' }}
                >
                    <X size={24} />
                </button>

                <div style={{ marginBottom: '1rem', paddingRight: '2rem' }}>
                    {isEditing ? (
                        <input
                            value={editTitle}
                            onChange={(e) => setEditTitle(e.target.value)}
                            style={{ fontSize: '1.8rem', width: '100%', padding: '0.5rem', marginBottom: '0.5rem' }}
                        />
                    ) : (
                        <h2 style={{ fontSize: '1.8rem', color: 'var(--primary-color)' }}>{note.title}</h2>
                    )}
                </div>

                <div style={{ display: 'flex', flexWrap: 'wrap', gap: '0.5rem', marginBottom: '1.5rem', alignItems: 'center' }}>
                    {note.tags && note.tags.map(tag => (
                        <span key={tag.tag_id} style={{
                            background: 'var(--bg-secondary)', color: 'var(--primary-color)',
                            padding: '0.25rem 0.75rem', borderRadius: '999px', fontSize: '0.85rem',
                            display: 'flex', alignItems: 'center', gap: '4px'
                        }}>
                            <TagIcon size={12} /> {tag.name}
                        </span>
                    ))}

                    <div style={{ position: 'relative' }}>
                        <button
                            onClick={() => setIsAddingTag(!isAddingTag)}
                            style={{
                                background: 'transparent', border: '1px dashed var(--border-color)',
                                borderRadius: '999px', padding: '0.25rem 0.75rem', fontSize: '0.85rem',
                                cursor: 'pointer', display: 'flex', alignItems: 'center', gap: '4px'
                            }}
                        >
                            <Plus size={14} /> Add Tag
                        </button>

                        {isAddingTag && (
                            <div style={{
                                position: 'absolute', top: '100%', left: 0, marginTop: '0.5rem',
                                background: 'white', border: '1px solid var(--border-color)', borderRadius: '8px',
                                boxShadow: 'var(--shadow-lg)', width: '200px', zIndex: 10,
                                maxHeight: '200px', overflowY: 'auto'
                            }}>
                                {availableTags.filter(t => !existingTagIds.includes(t.tag_id)).length === 0 ? (
                                    <div style={{ padding: '0.5rem', color: '#999', fontSize: '0.9rem' }}>No tags available</div>
                                ) : (
                                    availableTags.filter(t => !existingTagIds.includes(t.tag_id)).map(tag => (
                                        <div
                                            key={tag.tag_id}
                                            onClick={() => handleAddTag(tag.tag_id)}
                                            style={{
                                                padding: '0.5rem', cursor: 'pointer', fontSize: '0.9rem',
                                                borderBottom: '1px solid #f0f0f0'
                                            }}
                                            onMouseEnter={e => e.target.style.background = 'var(--bg-secondary)'}
                                            onMouseLeave={e => e.target.style.background = 'white'}
                                        >
                                            {tag.name}
                                        </div>
                                    ))
                                )}
                            </div>
                        )}
                    </div>
                </div>

                <div style={{ marginBottom: '2rem' }}>
                    {isEditing ? (
                        <textarea
                            value={editContent}
                            onChange={(e) => setEditContent(e.target.value)}
                            rows={10}
                            style={{ width: '100%', padding: '0.5rem', fontSize: '1rem', fontFamily: 'inherit' }}
                        />
                    ) : (
                        <div style={{ whiteSpace: 'pre-wrap', lineHeight: '1.6', color: '#333' }}>
                            {note.content}
                        </div>
                    )}
                </div>

                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderTop: '1px solid #eee', paddingTop: '1rem' }}>
                    <div style={{ fontSize: '0.8rem', color: '#999' }}>
                        Created: {new Date(note.created_at).toLocaleString()} | Last Updated: {new Date(note.updated_at).toLocaleString()}
                    </div>

                    <div style={{ display: 'flex', gap: '1rem' }}>
                        {isEditing ? (
                            <button
                                onClick={handleSave}
                                style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', padding: '0.5rem 1rem', background: 'var(--primary-color)', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer' }}
                            >
                                <Save size={16} /> Save
                            </button>
                        ) : (
                            <button
                                onClick={() => setIsEditing(true)}
                                style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', padding: '0.5rem 1rem', background: '#f0f0f0', color: '#333', border: 'none', borderRadius: '4px', cursor: 'pointer' }}
                            >
                                <Edit2 size={16} /> Edit
                            </button>
                        )}

                        <button
                            onClick={handleDelete}
                            style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', padding: '0.5rem 1rem', background: '#fee2e2', color: '#ef4444', border: 'none', borderRadius: '4px', cursor: 'pointer' }}
                        >
                            <Trash2 size={16} /> Delete
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default NoteDetailModal;
