import React, { useState, useEffect } from 'react';
import TenantForm from '../components/Tenants/TenantForm';
import TenantsList from '../components/Tenants/TenantsList';
import ConfirmDialog from '../components/common/ConfirmDialog';
import './Tenants.css';

// Import Go backend functions
import { 
  GetAllHouses,
  GetAllApartments,
  GetAllTenants,
  CreateTenant,
  UpdateTenant,
  DeleteTenant
} from '../../wailsjs/go/main/App';

const Tenants = () => {
  const [tenants, setTenants] = useState([]);
  const [houses, setHouses] = useState([]);
  const [apartments, setApartments] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isFormVisible, setIsFormVisible] = useState(false);
  const [currentTenant, setCurrentTenant] = useState(null);
  const [error, setError] = useState(null);
  const [infoMessage, setInfoMessage] = useState(null);
  const [confirmDialog, setConfirmDialog] = useState({
    isOpen: false,
    message: '',
    tenantId: null
  });
  
  // Fetch all data on component mount
  useEffect(() => {
    loadData();
  }, []);
  
  const loadData = async () => {
    try {
      setIsLoading(true);
      setError(null);
      
      // Load houses and apartments first
      const allHouses = await GetAllHouses();
      setHouses(allHouses || []);
      
      const allApartments = await GetAllApartments();
      setApartments(allApartments || []);
      
      // Check if we have both houses and apartments
      if ((allHouses && allHouses.length > 0) && (allApartments && allApartments.length > 0)) {
        const allTenants = await GetAllTenants();
        setTenants(allTenants || []);
      } else {
        setTenants([]);
        if (!allHouses || allHouses.length === 0) {
          setInfoMessage("Please add houses before adding tenants");
        } else if (!allApartments || allApartments.length === 0) {
          setInfoMessage("Please add apartments before adding tenants");
        }
      }
    } catch (err) {
      console.error('Error loading data:', err);
      setError('Failed to load data. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };
  
  const handleAddTenant = () => {
    if (!houses || houses.length === 0) {
      setError('Please add at least one house before adding tenants');
      return;
    }
    
    if (!apartments || apartments.length === 0) {
      setError('Please add at least one apartment before adding tenants');
      return;
    }
    
    setCurrentTenant(null);
    setIsFormVisible(true);
  };
  
  const handleEditTenant = (tenant) => {
    setCurrentTenant(tenant);
    setIsFormVisible(true);
  };
  
  const handleDeleteTenant = (id) => {
    setConfirmDialog({
      isOpen: true,
      message: 'Are you sure you want to delete this tenant?',
      tenantId: id
    });
  };
  
  const confirmDelete = async () => {
    try {
      setError(null);
      await DeleteTenant(confirmDialog.tenantId);
      setConfirmDialog({ isOpen: false, message: '', tenantId: null });
      await loadData();
    } catch (err) {
      console.error('Error deleting tenant:', err);
      setError('Failed to delete tenant. Please try again.');
      setConfirmDialog({ isOpen: false, message: '', tenantId: null });
    }
  };
  
  const cancelDelete = () => {
    setConfirmDialog({ isOpen: false, message: '', tenantId: null });
  };
  
  const handleSaveTenant = async (tenantData) => {
    try {
      setError(null);
      
      if (tenantData.id) {
        // Update existing tenant
        await UpdateTenant(
          tenantData.id,
          tenantData.firstName,
          tenantData.lastName,
          tenantData.moveInDate,
          tenantData.moveOutDate,
          tenantData.deposit,
          tenantData.email,
          tenantData.numberOfPersons,
          tenantData.targetColdRent,
          tenantData.targetAncillaryPayment,
          tenantData.targetElectricityPayment,
          tenantData.greeting,
          tenantData.houseId,
          tenantData.apartmentId
        );
      } else {
        // Create new tenant
        await CreateTenant(
          tenantData.firstName,
          tenantData.lastName,
          tenantData.moveInDate,
          tenantData.moveOutDate,
          tenantData.deposit,
          tenantData.email,
          tenantData.numberOfPersons,
          tenantData.targetColdRent,
          tenantData.targetAncillaryPayment,
          tenantData.targetElectricityPayment,
          tenantData.greeting,
          tenantData.houseId,
          tenantData.apartmentId
        );
      }
      
      // Reload tenants and close form
      await loadData();
      setIsFormVisible(false);
    } catch (err) {
      console.error('Error saving tenant:', err);
      setError(`Failed to ${tenantData.id ? 'update' : 'add'} tenant: ${err.message || err}`);
    }
  };
  
  const handleCancelForm = () => {
    setIsFormVisible(false);
    setCurrentTenant(null);
  };
  
  return (
    <div className="tenants-page">
      <div className="tenants-header">
        <h1>Manage Tenants</h1>
        {!isFormVisible && (
          <button className="add-button" onClick={handleAddTenant}>
            Add New Tenant
          </button>
        )}
      </div>
      
      {error && <div className="error-banner">{error}</div>}
      {infoMessage && <div className="info-banner">{infoMessage}</div>}
      
      {isLoading ? (
        <div className="loading">Loading tenants...</div>
      ) : (
        <div className="tenants-content">
          {isFormVisible ? (
            <TenantForm 
              tenant={currentTenant} 
              houses={houses}
              apartments={apartments}
              onSave={handleSaveTenant} 
              onCancel={handleCancelForm} 
            />
          ) : (
            <TenantsList 
              tenants={tenants} 
              houses={houses}
              apartments={apartments}
              onEdit={handleEditTenant} 
              onDelete={handleDeleteTenant} 
            />
          )}
        </div>
      )}
      
      <ConfirmDialog 
        isOpen={confirmDialog.isOpen}
        message={confirmDialog.message}
        onConfirm={confirmDelete}
        onCancel={cancelDelete}
      />
    </div>
  );
};

export default Tenants;