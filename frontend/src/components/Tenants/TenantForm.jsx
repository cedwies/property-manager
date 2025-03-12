import React, { useState, useEffect } from 'react';
import './TenantForm.css';

const TenantForm = ({ tenant, houses, apartments, onSave, onCancel }) => {
  const [formData, setFormData] = useState({
    id: null,
    firstName: '',
    lastName: '',
    moveInDate: '',
    moveOutDate: '',
    deposit: '',
    email: '',
    numberOfPersons: '1',
    targetColdRent: '',
    targetAncillaryPayment: '',
    targetElectricityPayment: '',
    greeting: '',
    houseId: '',
    apartmentId: '',
  });
  
  const [errors, setErrors] = useState({});
  const [housesApartments, setHousesApartments] = useState({});
  const [filteredApartments, setFilteredApartments] = useState([]);
  
  // Initialize form with tenant data if provided (for editing)
  useEffect(() => {
    if (tenant) {
      // Format dates for form input
      const moveInDate = tenant.moveInDate ? new Date(tenant.moveInDate).toISOString().split('T')[0] : '';
      const moveOutDate = tenant.moveOutDate ? new Date(tenant.moveOutDate).toISOString().split('T')[0] : '';
      
      setFormData({
        id: tenant.id,
        firstName: tenant.firstName || '',
        lastName: tenant.lastName || '',
        moveInDate: moveInDate,
        moveOutDate: moveOutDate,
        deposit: tenant.deposit ? tenant.deposit.toString() : '',
        email: tenant.email || '',
        numberOfPersons: tenant.numberOfPersons ? tenant.numberOfPersons.toString() : '1',
        targetColdRent: tenant.targetColdRent ? tenant.targetColdRent.toString() : '',
        targetAncillaryPayment: tenant.targetAncillaryPayment ? tenant.targetAncillaryPayment.toString() : '',
        targetElectricityPayment: tenant.targetElectricityPayment ? tenant.targetElectricityPayment.toString() : '',
        greeting: tenant.greeting || '',
        houseId: tenant.houseId ? tenant.houseId.toString() : '',
        apartmentId: tenant.apartmentId ? tenant.apartmentId.toString() : '',
      });
    } else if (houses && houses.length > 0) {
      // Set default house ID for new tenants if houses are available
      setFormData(prev => ({
        ...prev,
        houseId: houses[0].id.toString()
      }));
    }
  }, [tenant, houses]);
  
  // Organize apartments by house ID
  useEffect(() => {
    if (apartments && apartments.length > 0) {
      const apartmentsByHouse = {};
      
      apartments.forEach(apartment => {
        const houseId = apartment.houseId.toString();
        if (!apartmentsByHouse[houseId]) {
          apartmentsByHouse[houseId] = [];
        }
        apartmentsByHouse[houseId].push(apartment);
      });
      
      setHousesApartments(apartmentsByHouse);
    }
  }, [apartments]);
  
  // Filter apartments when house changes
  useEffect(() => {
    if (formData.houseId && housesApartments[formData.houseId]) {
      setFilteredApartments(housesApartments[formData.houseId]);
      
      // If current apartment doesn't belong to selected house, reset it
      const apartmentBelongsToHouse = housesApartments[formData.houseId].some(
        apt => apt.id.toString() === formData.apartmentId
      );
      
      if (!apartmentBelongsToHouse && housesApartments[formData.houseId].length > 0) {
        setFormData(prev => ({
          ...prev,
          apartmentId: housesApartments[formData.houseId][0].id.toString()
        }));
      }
    } else {
      setFilteredApartments([]);
      setFormData(prev => ({
        ...prev,
        apartmentId: ''
      }));
    }
  }, [formData.houseId, housesApartments]);
  
  const validate = () => {
    const newErrors = {};
    
    // Basic validations
    if (!formData.firstName.trim()) {
      newErrors.firstName = 'First name is required';
    }
    
    if (!formData.lastName.trim()) {
      newErrors.lastName = 'Last name is required';
    }
    
    if (!formData.moveInDate) {
      newErrors.moveInDate = 'Move-in date is required';
    }
    
    // If move-out date is provided, ensure it's after move-in date
    if (formData.moveOutDate) {
      const moveIn = new Date(formData.moveInDate);
      const moveOut = new Date(formData.moveOutDate);
      
      if (moveOut <= moveIn) {
        newErrors.moveOutDate = 'Move-out date must be after move-in date';
      }
    }
    
    // Numeric validations
    if (formData.deposit.trim() !== '') {
      const deposit = parseFloat(formData.deposit.replace(',', '.'));
      if (isNaN(deposit) || deposit < 0) {
        newErrors.deposit = 'Deposit must be a positive number';
      }
    }
    
    if (!formData.numberOfPersons.trim()) {
      newErrors.numberOfPersons = 'Number of persons is required';
    } else {
      const persons = parseInt(formData.numberOfPersons);
      if (isNaN(persons) || persons <= 0) {
        newErrors.numberOfPersons = 'Number of persons must be greater than 0';
      }
    }
    
    if (!formData.targetColdRent.trim()) {
      newErrors.targetColdRent = 'Target cold rent is required';
    } else {
      const rent = parseFloat(formData.targetColdRent.replace(',', '.'));
      if (isNaN(rent) || rent <= 0) {
        newErrors.targetColdRent = 'Target cold rent must be greater than 0';
      }
    }
    
    if (!formData.targetAncillaryPayment.trim()) {
      newErrors.targetAncillaryPayment = 'Target ancillary payment is required';
    } else {
      const payment = parseFloat(formData.targetAncillaryPayment.replace(',', '.'));
      if (isNaN(payment) || payment <= 0) {
        newErrors.targetAncillaryPayment = 'Target ancillary payment must be greater than 0';
      }
    }
    
    if (!formData.targetElectricityPayment.trim()) {
      newErrors.targetElectricityPayment = 'Target electricity payment is required';
    } else {
      const payment = parseFloat(formData.targetElectricityPayment.replace(',', '.'));
      if (isNaN(payment) || payment <= 0) {
        newErrors.targetElectricityPayment = 'Target electricity payment must be greater than 0';
      }
    }
    
    // House and apartment validations
    if (!formData.houseId) {
      newErrors.houseId = 'Please select a house';
    }
    
    if (!formData.apartmentId) {
      newErrors.apartmentId = 'Please select an apartment';
    }
    
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };
  
  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value,
    });
    
    // Clear error when field is edited
    if (errors[name]) {
      setErrors({
        ...errors,
        [name]: null,
      });
    }
    
    // Special handling for move-in date to update move-out date minimum
    if (name === 'moveInDate' && value && formData.moveOutDate) {
      // Check if move-out date is before move-in date
      const moveIn = new Date(value);
      const moveOut = new Date(formData.moveOutDate);
      
      if (moveOut <= moveIn) {
        // Clear move-out date if it's before the new move-in date
        setFormData(prev => ({
          ...prev,
          moveOutDate: ''
        }));
      }
    }
  };
  
  const handleSubmit = (e) => {
    e.preventDefault();
    
    if (validate()) {
      onSave({
        ...formData,
        houseId: parseInt(formData.houseId),
        apartmentId: parseInt(formData.apartmentId),
      });
    }
  };
  
  return (
    <form className="tenant-form" onSubmit={handleSubmit}>
      <h2>{formData.id ? 'Edit Tenant' : 'Add New Tenant'}</h2>
      
      <div className="form-section">
        <h3>Basic Information</h3>
        <div className="form-row">
          <div className="form-group medium">
            <label htmlFor="firstName">First Name*</label>
            <input
              type="text"
              id="firstName"
              name="firstName"
              value={formData.firstName}
              onChange={handleChange}
              className={errors.firstName ? 'error' : ''}
              placeholder="John"
            />
            {errors.firstName && <div className="error-message">{errors.firstName}</div>}
          </div>
          
          <div className="form-group medium">
            <label htmlFor="lastName">Last Name*</label>
            <input
              type="text"
              id="lastName"
              name="lastName"
              value={formData.lastName}
              onChange={handleChange}
              className={errors.lastName ? 'error' : ''}
              placeholder="Doe"
            />
            {errors.lastName && <div className="error-message">{errors.lastName}</div>}
          </div>
        </div>
        
        <div className="form-row">
          <div className="form-group medium">
            <label htmlFor="email">Email (Optional)</label>
            <input
              type="email"
              id="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              className={errors.email ? 'error' : ''}
              placeholder="john.doe@example.com"
            />
            {errors.email && <div className="error-message">{errors.email}</div>}
          </div>
          
          <div className="form-group medium">
            <label htmlFor="numberOfPersons">Number of Persons*</label>
            <input
              type="number"
              id="numberOfPersons"
              name="numberOfPersons"
              value={formData.numberOfPersons}
              onChange={handleChange}
              min="1"
              className={errors.numberOfPersons ? 'error' : ''}
            />
            {errors.numberOfPersons && <div className="error-message">{errors.numberOfPersons}</div>}
          </div>
        </div>
        
        <div className="form-group">
          <label htmlFor="greeting">Greeting/Salutation*</label>
          <input
            type="text"
            id="greeting"
            name="greeting"
            value={formData.greeting}
            onChange={handleChange}
            className={errors.greeting ? 'error' : ''}
            placeholder="e.g., Dear Mr. Doe"
          />
          {errors.greeting && <div className="error-message">{errors.greeting}</div>}
        </div>
      </div>
      
      <div className="form-section">
        <h3>Location</h3>
        <div className="form-group">
          <label htmlFor="houseId">House*</label>
          <select
            id="houseId"
            name="houseId"
            value={formData.houseId}
            onChange={handleChange}
            className={errors.houseId ? 'error' : ''}
          >
            <option value="">Select a house</option>
            {houses && houses.map(house => (
              <option key={house.id} value={house.id.toString()}>
                {house.name} ({house.street} {house.number}, {house.city})
              </option>
            ))}
          </select>
          {errors.houseId && <div className="error-message">{errors.houseId}</div>}
        </div>
        
        <div className="form-group">
          <label htmlFor="apartmentId">Apartment*</label>
          <select
            id="apartmentId"
            name="apartmentId"
            value={formData.apartmentId}
            onChange={handleChange}
            className={errors.apartmentId ? 'error' : ''}
            disabled={!formData.houseId || filteredApartments.length === 0}
          >
            <option value="">Select an apartment</option>
            {filteredApartments.map(apartment => (
              <option key={apartment.id} value={apartment.id.toString()}>
                {apartment.name} ({apartment.size} m²)
              </option>
            ))}
          </select>
          {errors.apartmentId && <div className="error-message">{errors.apartmentId}</div>}
          {formData.houseId && filteredApartments.length === 0 && 
            <div className="error-message">No apartments available for this house</div>
          }
        </div>
      </div>
      
      <div className="form-section">
        <h3>Dates and Financial Details</h3>
        <div className="form-row">
          <div className="form-group medium">
            <label htmlFor="moveInDate">
              Move-in Date*
              <span className="date-tooltip">ⓘ
                <span className="tooltip-text">Click to open a calendar for selecting the date when the tenant moved in</span>
              </span>
            </label>
            <input
              type="date"
              id="moveInDate"
              name="moveInDate"
              value={formData.moveInDate}
              onChange={handleChange}
              className={errors.moveInDate ? 'error' : ''}
              data-date-format="YYYY-MM-DD"
            />
            {errors.moveInDate && <div className="error-message">{errors.moveInDate}</div>}
          </div>
          
          <div className="form-group medium">
            <label htmlFor="moveOutDate">
              Move-out Date (Optional)
              <span className="date-tooltip">ⓘ
                <span className="tooltip-text">Click to select the date when the tenant moved/will move out. Leave empty for current tenants.</span>
              </span>
            </label>
            <div className="date-input-container">
              <input
                type="date"
                id="moveOutDate"
                name="moveOutDate"
                value={formData.moveOutDate}
                onChange={handleChange}
                className={errors.moveOutDate ? 'error' : ''}
                data-date-format="YYYY-MM-DD"
              />
              {formData.moveOutDate && (
                <button 
                  type="button"
                  className="clear-date-btn"
                  onClick={() => {
                    setFormData({
                      ...formData,
                      moveOutDate: ''
                    });
                    
                    // Clear any errors
                    if (errors.moveOutDate) {
                      setErrors({
                        ...errors,
                        moveOutDate: null
                      });
                    }
                  }}
                  title="Clear move-out date"
                >
                  ×
                </button>
              )}
            </div>
            {errors.moveOutDate && <div className="error-message">{errors.moveOutDate}</div>}
            <div className="helper-text">Leave empty for current tenants</div>
          </div>
        </div>
        
        <div className="form-group">
          <label htmlFor="deposit">Deposit*</label>
          <input
            type="text"
            id="deposit"
            name="deposit"
            value={formData.deposit}
            onChange={handleChange}
            className={errors.deposit ? 'error' : ''}
            placeholder="e.g., 1500 or 1500,00"
          />
          {errors.deposit && <div className="error-message">{errors.deposit}</div>}
        </div>
        
        <div className="form-row">
          <div className="form-group">
            <label htmlFor="targetColdRent">Target Cold Rent*</label>
            <input
              type="text"
              id="targetColdRent"
              name="targetColdRent"
              value={formData.targetColdRent}
              onChange={handleChange}
              className={errors.targetColdRent ? 'error' : ''}
              placeholder="e.g., 800 or 800,50"
            />
            {errors.targetColdRent && <div className="error-message">{errors.targetColdRent}</div>}
          </div>
          
          <div className="form-group">
            <label htmlFor="targetAncillaryPayment">Target Ancillary Payment*</label>
            <input
              type="text"
              id="targetAncillaryPayment"
              name="targetAncillaryPayment"
              value={formData.targetAncillaryPayment}
              onChange={handleChange}
              className={errors.targetAncillaryPayment ? 'error' : ''}
              placeholder="e.g., 200 or 200,00"
            />
            {errors.targetAncillaryPayment && <div className="error-message">{errors.targetAncillaryPayment}</div>}
          </div>
          
          <div className="form-group">
            <label htmlFor="targetElectricityPayment">Target Electricity Payment*</label>
            <input
              type="text"
              id="targetElectricityPayment"
              name="targetElectricityPayment"
              value={formData.targetElectricityPayment}
              onChange={handleChange}
              className={errors.targetElectricityPayment ? 'error' : ''}
              placeholder="e.g., 100 or 100,00"
            />
            {errors.targetElectricityPayment && <div className="error-message">{errors.targetElectricityPayment}</div>}
          </div>
        </div>
      </div>
      
      <div className="form-buttons">
        <button type="button" className="button secondary" onClick={onCancel}>
          Cancel
        </button>
        <button type="submit" className="button primary">
          {formData.id ? 'Update' : 'Add'} Tenant
        </button>
      </div>
    </form>
  );
};

export default TenantForm;