import React, { useState, useEffect } from 'react';
import './PaymentRecordRow.css';

const PaymentRecordRow = ({ record, onChange, onOpenNote, onToggleLock }) => {
  const [localRecord, setLocalRecord] = useState({ ...record });
  const [isNewRecord, setIsNewRecord] = useState(record.id === 0);

  useEffect(() => {
    setLocalRecord({ ...record });
    setIsNewRecord(record.id === 0);
  }, [record]);

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    
    // Don't allow changes if record is locked
    if (localRecord.isLocked) return;
    
    setLocalRecord({
      ...localRecord,
      [name]: value
    });
  };

  const handleBlur = () => {
    // Always send changes if something changed, even for new records
    if (JSON.stringify(localRecord) !== JSON.stringify(record)) {
      // Create a copy of the record with numeric values for payment fields
      const processedRecord = {
        ...localRecord,
        paidColdRent: parseFloat(localRecord.paidColdRent) || 0,
        paidAncillary: parseFloat(localRecord.paidAncillary) || 0,
        paidElectricity: parseFloat(localRecord.paidElectricity) || 0,
        extraPayments: parseFloat(localRecord.extraPayments) || 0,
        persons: parseInt(localRecord.persons) || 0
      };
      
      onChange(processedRecord);
    }
  };

  const formatMonth = (monthStr) => {
    // Convert YYYY-MM to MM.YYYY format
    const [year, month] = monthStr.split('-');
    return `${month}.${year}`;
  };

  const handleNoteClick = () => {
    onOpenNote(localRecord.id, localRecord.note);
  };

  const handleLockClick = () => {
    // Create a copy with numeric values for payment fields
    const processedRecord = {
      ...localRecord,
      paidColdRent: parseFloat(localRecord.paidColdRent) || 0,
      paidAncillary: parseFloat(localRecord.paidAncillary) || 0,
      paidElectricity: parseFloat(localRecord.paidElectricity) || 0,
      extraPayments: parseFloat(localRecord.extraPayments) || 0,
      persons: parseInt(localRecord.persons) || 0
    };
    
    onToggleLock(processedRecord.id);
  };

  // Get current month for highlighting
  const now = new Date();
  const currentMonth = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`;
  const isCurrentMonth = localRecord.month === currentMonth;

  // Determine if this is a future month
  const isFutureMonth = localRecord.month > currentMonth;

  return (
    <div className={`payment-record-row ${localRecord.isLocked ? 'locked' : ''} ${isCurrentMonth ? 'current-month' : ''}`}>
      <div className="row-cell month-cell">{formatMonth(localRecord.month)}</div>
      
      <div className="row-cell">€{localRecord.targetColdRent.toFixed(2)}</div>
      
      <div className="row-cell">
        <input
          type="text"
          name="paidColdRent"
          value={localRecord.paidColdRent === 0 && !localRecord.isLocked ? '' : localRecord.paidColdRent}
          onChange={handleInputChange}
          onBlur={handleBlur}
          disabled={localRecord.isLocked || isFutureMonth}
          placeholder="0.00"
          className={`amount-input ${localRecord.isLocked || isFutureMonth ? 'disabled' : ''}`}
        />
      </div>
      
      <div className="row-cell">
        <input
          type="text"
          name="paidAncillary"
          value={localRecord.paidAncillary === 0 && !localRecord.isLocked ? '' : localRecord.paidAncillary}
          onChange={handleInputChange}
          onBlur={handleBlur}
          disabled={localRecord.isLocked || isFutureMonth}
          placeholder="0.00"
          className={`amount-input ${localRecord.isLocked || isFutureMonth ? 'disabled' : ''}`}
        />
      </div>
      
      <div className="row-cell">
        <input
          type="text"
          name="paidElectricity"
          value={localRecord.paidElectricity === 0 && !localRecord.isLocked ? '' : localRecord.paidElectricity}
          onChange={handleInputChange}
          onBlur={handleBlur}
          disabled={localRecord.isLocked || isFutureMonth}
          placeholder="0.00"
          className={`amount-input ${localRecord.isLocked || isFutureMonth ? 'disabled' : ''}`}
        />
      </div>
      
      <div className="row-cell">
        <input
          type="text"
          name="extraPayments"
          value={localRecord.extraPayments === 0 && !localRecord.isLocked ? '' : localRecord.extraPayments}
          onChange={handleInputChange}
          onBlur={handleBlur}
          disabled={localRecord.isLocked || isFutureMonth}
          placeholder="0.00"
          className={`amount-input ${localRecord.isLocked || isFutureMonth ? 'disabled' : ''}`}
        />
      </div>
      
      <div className="row-cell">
        <input
          type="number"
          name="persons"
          value={localRecord.persons || ''}
          onChange={handleInputChange}
          onBlur={handleBlur}
          disabled={localRecord.isLocked || isFutureMonth}
          min="1"
          className={`persons-input ${localRecord.isLocked || isFutureMonth ? 'disabled' : ''}`}
        />
      </div>
      
      <div className="row-cell action-cell">
        <button 
          className={`note-button ${localRecord.note ? 'has-note' : ''} ${isFutureMonth ? 'disabled' : ''}`}
          onClick={handleNoteClick}
          disabled={isFutureMonth}
          title={localRecord.note ? 'View/Edit Note' : 'Add Note'}
        >
          {localRecord.note ? '📝' : '✏️'}
        </button>
      </div>
      
      <div className="row-cell action-cell">
        <button 
          className={`lock-button ${localRecord.isLocked ? 'locked' : ''} ${isFutureMonth ? 'disabled' : ''}`}
          onClick={handleLockClick}
          disabled={isFutureMonth}
          title={localRecord.isLocked ? 'Unlock Record' : 'Lock Record'}
        >
          {localRecord.isLocked ? '🔒' : '🔓'}
        </button>
      </div>
    </div>
  );
};

export default PaymentRecordRow;