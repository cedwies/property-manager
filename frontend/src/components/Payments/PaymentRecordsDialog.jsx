import React, { useState, useEffect } from 'react';
import './PaymentRecordsDialog.css';
import PaymentRecordsList from './PaymentRecordsList';
import NoteDialog from './NoteDialog';

import { 
  GetCurrentTenantsByHouseID, 
  GetPaymentRecordsForHouse,
  BatchSavePaymentRecords,
  GetLastTwelveMonths
} from '../../../wailsjs/go/main/App';

const PaymentRecordsDialog = ({ house, onClose }) => {
  const [tenants, setTenants] = useState([]);
  const [paymentData, setPaymentData] = useState({});
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState(null);
  const [currentTenantIndex, setCurrentTenantIndex] = useState(0);
  const [months, setMonths] = useState([]);
  const [showNoteDialog, setShowNoteDialog] = useState(false);
  const [currentNote, setCurrentNote] = useState('');
  const [currentRecordId, setCurrentRecordId] = useState(null);
  const [hasChanges, setHasChanges] = useState(false);

  // Load tenants and payment data
  useEffect(() => {
    loadData();
  }, [house.id]);

  const loadData = async () => {
    try {
      setIsLoading(true);
      setError(null);

      // Get last twelve months
      const monthsList = await GetLastTwelveMonths();
      setMonths(monthsList);

      // Get current tenants for this house
      const currentTenants = await GetCurrentTenantsByHouseID(house.id);
      setTenants(currentTenants || []);

      // Get payment records for house tenants
      const payments = await GetPaymentRecordsForHouse(house.id);
      setPaymentData(payments || {});
      
      setHasChanges(false);
      setCurrentTenantIndex(0);
    } catch (err) {
      console.error('Error loading payment data:', err);
      setError('Failed to load payment data. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleNextTenant = async () => {
    // Save current tenant's changes first
    if (hasChanges) {
      await saveChanges();
    }
    
    // Move to next tenant
    if (currentTenantIndex < tenants.length - 1) {
      setCurrentTenantIndex(currentTenantIndex + 1);
      setHasChanges(false);
    } else {
      // All tenants checked, close dialog
      onClose();
      alert(`All tenants for ${house.name} have been checked.`);
    }
  };

  const handleConfirm = async () => {
    // Save changes and close
    if (hasChanges) {
      await saveChanges();
    }
    onClose();
  };

  const handleCancel = () => {
    // Confirm discard changes if there are any
    if (hasChanges) {
      if (window.confirm('You have unsaved changes. Are you sure you want to discard them?')) {
        onClose();
      }
    } else {
      onClose();
    }
  };

  // Helper function to ensure numeric values
  const ensureNumericValues = (record) => {
    return {
      ...record,
      paidColdRent: typeof record.paidColdRent === 'string' ? 
        parseFloat(record.paidColdRent) || 0 : record.paidColdRent,
      paidAncillary: typeof record.paidAncillary === 'string' ? 
        parseFloat(record.paidAncillary) || 0 : record.paidAncillary,
      paidElectricity: typeof record.paidElectricity === 'string' ? 
        parseFloat(record.paidElectricity) || 0 : record.paidElectricity,
      extraPayments: typeof record.extraPayments === 'string' ? 
        parseFloat(record.extraPayments) || 0 : record.extraPayments,
      persons: typeof record.persons === 'string' ? 
        parseInt(record.persons) || 0 : record.persons
    };
  };

  const saveChanges = async () => {
    try {
      setIsSaving(true);
      
      // Create an array of payment records to save
      const currentTenant = tenants[currentTenantIndex];
      if (!currentTenant) return;
      
      const records = paymentData[currentTenant.id] || [];
      if (records.length === 0) return;
      
      // Filter out records with ID 0 that have no actual data entered
      const recordsToSave = records
        .filter(record => {
          // Include all existing records (ID > 0)
          if (record.id > 0) return true;
          
          // For new records (ID = 0), only include those with data
          return (
            record.paidColdRent > 0 || 
            record.paidAncillary > 0 || 
            record.paidElectricity > 0 || 
            record.extraPayments > 0 || 
            record.note.trim() !== '' ||
            record.isLocked
          );
        })
        .map(record => ensureNumericValues(record)); // Ensure numeric values
      
      if (recordsToSave.length > 0) {
        await BatchSavePaymentRecords(recordsToSave);
      }
      
      setHasChanges(false);
    } catch (err) {
      console.error('Error saving payment records:', err);
      setError('Failed to save payment records. Please try again.');
    } finally {
      setIsSaving(false);
    }
  };

  const handleRecordChange = (updatedRecord) => {
    setHasChanges(true);
    
    // Update the record in the payment data
    const currentTenant = tenants[currentTenantIndex];
    if (!currentTenant) return;
    
    const records = [...(paymentData[currentTenant.id] || [])];
    const recordIndex = records.findIndex(r => r.id === updatedRecord.id && r.month === updatedRecord.month);
    
    if (recordIndex !== -1) {
      // Update existing record
      records[recordIndex] = updatedRecord;
    } else {
      // Add new record
      records.push(updatedRecord);
    }
    
    // Update the payment data
    setPaymentData({
      ...paymentData,
      [currentTenant.id]: records
    });
  };

  const handleOpenNoteDialog = (recordId, note) => {
    setCurrentRecordId(recordId);
    setCurrentNote(note);
    setShowNoteDialog(true);
  };

  const handleSaveNote = (note) => {
    // Find and update the record
    const currentTenant = tenants[currentTenantIndex];
    if (!currentTenant) return;
    
    const records = [...(paymentData[currentTenant.id] || [])];
    const recordIndex = records.findIndex(r => r.id === currentRecordId || 
      (r.id === 0 && currentRecordId === 0 && r.month === records.find(rec => rec.id === 0)?.month));
    
    if (recordIndex !== -1) {
      records[recordIndex] = {
        ...records[recordIndex],
        note: note
      };
      
      // Update the payment data
      setPaymentData({
        ...paymentData,
        [currentTenant.id]: records
      });
      
      setHasChanges(true);
    }
    
    setShowNoteDialog(false);
  };

  const handleToggleLock = (recordId) => {
    // Find and update the record's lock status
    const currentTenant = tenants[currentTenantIndex];
    if (!currentTenant) return;
    
    const records = [...(paymentData[currentTenant.id] || [])];
    const recordIndex = records.findIndex(r => r.id === recordId || 
      (r.id === 0 && recordId === 0 && r.month === records.find(rec => rec.id === 0)?.month));
    
    if (recordIndex !== -1) {
      records[recordIndex] = {
        ...records[recordIndex],
        isLocked: !records[recordIndex].isLocked
      };
      
      // Update the payment data
      setPaymentData({
        ...paymentData,
        [currentTenant.id]: records
      });
      
      setHasChanges(true);
    }
  };

  const currentTenant = tenants[currentTenantIndex] || null;
  const currentRecords = currentTenant ? (paymentData[currentTenant.id] || []) : [];

  return (
    <div className="payment-dialog-overlay">
      <div className="payment-dialog">
        <div className="payment-dialog-header">
          <h2>Payment Records - {house.name}</h2>
          <button 
            className="close-button"
            onClick={handleCancel}
            title="Close dialog"
          >
            ×
          </button>
        </div>

        {error && <div className="error-banner">{error}</div>}

        {isLoading ? (
          <div className="loading">Loading payment records...</div>
        ) : isSaving ? (
          <div className="loading">Saving payment records...</div>
        ) : (
          <>
            {tenants.length === 0 ? (
              <div className="no-tenants-message">
                <p>No current tenants found for this house.</p>
              </div>
            ) : (
              <div className="payment-dialog-content">
                <div className="tenant-info">
                  <h3>
                    Tenant: {currentTenant?.firstName} {currentTenant?.lastName}
                    <span className="tenant-counter">
                      ({currentTenantIndex + 1} of {tenants.length})
                    </span>
                  </h3>
                  <div className="tenant-details">
                    <div>
                      <strong>Apartment:</strong> {currentTenant?.apartment?.name}
                    </div>
                    <div>
                      <strong>Target Cold Rent:</strong> €{currentTenant?.targetColdRent.toFixed(2)}
                    </div>
                    <div>
                      <strong>Persons:</strong> {currentTenant?.numberOfPersons}
                    </div>
                  </div>
                </div>

                <PaymentRecordsList 
                  records={currentRecords}
                  months={months}
                  onRecordChange={handleRecordChange}
                  onOpenNote={handleOpenNoteDialog}
                  onToggleLock={handleToggleLock}
                  tenantMoveInDate={currentTenant?.moveInDate}
                  currentTenant={currentTenant}
                />

                <div className="payment-dialog-actions">
                  <button 
                    className="button secondary"
                    onClick={handleCancel}
                  >
                    Cancel
                  </button>
                  <button 
                    className="button primary" 
                    onClick={handleConfirm}
                  >
                    Confirm
                  </button>
                  {currentTenantIndex < tenants.length - 1 && (
                    <button 
                      className="button next-tenant" 
                      onClick={handleNextTenant}
                    >
                      Next Tenant
                    </button>
                  )}
                  {currentTenantIndex === tenants.length - 1 && (
                    <button 
                      className="button next-tenant" 
                      onClick={handleNextTenant}
                    >
                      Finish All
                    </button>
                  )}
                </div>
              </div>
            )}
          </>
        )}
      </div>

      {showNoteDialog && (
        <NoteDialog 
          note={currentNote}
          onSave={handleSaveNote}
          onCancel={() => setShowNoteDialog(false)}
        />
      )}
    </div>
  );
};

export default PaymentRecordsDialog;