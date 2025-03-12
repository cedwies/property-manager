import React, { useState, useEffect } from 'react';
import './TenantsList.css';

const TenantsList = ({ tenants, houses, apartments, onEdit, onDelete }) => {
  const [filteredTenants, setFilteredTenants] = useState([]);
  const [selectedHouseId, setSelectedHouseId] = useState('all');
  const [selectedApartmentId, setSelectedApartmentId] = useState('all');
  const [availableApartments, setAvailableApartments] = useState([]);

  useEffect(() => {
    if (tenants && tenants.length > 0) {
      filterTenants();
    } else {
      setFilteredTenants([]);
    }
  }, [tenants, selectedHouseId, selectedApartmentId]);

  useEffect(() => {
    // When house selection changes, update available apartments
    if (selectedHouseId === 'all') {
      setAvailableApartments(apartments || []);
    } else {
      const houseApartments = apartments ? apartments.filter(
        apartment => apartment.houseId.toString() === selectedHouseId
      ) : [];
      setAvailableApartments(houseApartments);
    }

    // Reset apartment selection if the current selection doesn't belong to selected house
    if (selectedHouseId !== 'all' && selectedApartmentId !== 'all') {
      const apartmentBelongsToHouse = apartments ? apartments.some(
        apartment => apartment.id.toString() === selectedApartmentId && 
                    apartment.houseId.toString() === selectedHouseId
      ) : false;
      
      if (!apartmentBelongsToHouse) {
        setSelectedApartmentId('all');
      }
    }
  }, [selectedHouseId, apartments]);

  const filterTenants = () => {
    if (!tenants) return [];

    let filtered = [...tenants];
    
    // Filter by house
    if (selectedHouseId !== 'all') {
      filtered = filtered.filter(tenant => 
        tenant.houseId.toString() === selectedHouseId
      );
    }
    
    // Filter by apartment
    if (selectedApartmentId !== 'all') {
      filtered = filtered.filter(tenant => 
        tenant.apartmentId.toString() === selectedApartmentId
      );
    }
    
    setFilteredTenants(filtered);
  };

  const handleHouseFilterChange = (e) => {
    setSelectedHouseId(e.target.value);
  };

  const handleApartmentFilterChange = (e) => {
    setSelectedApartmentId(e.target.value);
  };

  const formatDate = (dateString) => {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString();
  };

  const formatCurrency = (amount) => {
    return amount ? `€${amount.toFixed(2).replace('.', ',')}` : '';
  };

  const getTenantStatus = (tenant) => {
    const now = new Date();
    
    if (tenant.moveOutDate && new Date(tenant.moveOutDate) < now) {
      return <span className="status-former">Former</span>;
    }
    
    return <span className="status-active">Active</span>;
  };

  if (!tenants || tenants.length === 0) {
    return (
      <div className="tenants-list empty">
        <p>No tenants found. Add your first tenant using the form.</p>
      </div>
    );
  }

  return (
    <div className="tenants-list">
      <div className="tenants-header">
        <h2>Your Tenants</h2>
        <div className="filter-container">
          <label htmlFor="house-filter">House:</label>
          <select
            id="house-filter"
            value={selectedHouseId}
            onChange={handleHouseFilterChange}
          >
            <option value="all">All Houses</option>
            {houses && houses.map(house => (
              <option key={house.id} value={house.id.toString()}>
                {house.name}
              </option>
            ))}
          </select>
          
          <label htmlFor="apartment-filter">Apartment:</label>
          <select
            id="apartment-filter"
            value={selectedApartmentId}
            onChange={handleApartmentFilterChange}
            disabled={selectedHouseId === 'all' && availableApartments.length === 0}
          >
            <option value="all">All Apartments</option>
            {availableApartments.map(apartment => (
              <option key={apartment.id} value={apartment.id.toString()}>
                {apartment.name}
              </option>
            ))}
          </select>
        </div>
      </div>

      <div className="list-container">
        <table>
          <thead>
            <tr>
              <th>Tenant</th>
              <th>Status</th>
              <th>Location</th>
              <th>Move-in</th>
              <th>Move-out</th>
              <th>Persons</th>
              <th>Cold Rent</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {filteredTenants.map((tenant) => (
              <tr key={tenant.id}>
                <td>
                  <div className="tenant-name">{tenant.firstName} {tenant.lastName}</div>
                  {tenant.email && <div className="tenant-email">{tenant.email}</div>}
                </td>
                <td>{getTenantStatus(tenant)}</td>
                <td>
                  {tenant.house?.name}, {tenant.apartment?.name}
                </td>
                <td className="date-cell">{formatDate(tenant.moveInDate)}</td>
                <td className="date-cell">{formatDate(tenant.moveOutDate)}</td>
                <td>{tenant.numberOfPersons}</td>
                <td className="payment-cell">{formatCurrency(tenant.targetColdRent)}</td>
                <td className="actions">
                  <button 
                    className="button edit" 
                    onClick={() => onEdit(tenant)}
                    title="Edit tenant"
                  >
                    Edit
                  </button>
                  <button 
                    className="button delete" 
                    onClick={() => onDelete(tenant.id)}
                    title="Delete tenant"
                  >
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default TenantsList;