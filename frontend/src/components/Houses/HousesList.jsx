import React from 'react';
import './HousesList.css';

const HousesList = ({ houses, onEdit, onDelete }) => {
  if (!houses || houses.length === 0) {
    return (
      <div className="houses-list empty">
        <p>No houses found. Add your first house using the form.</p>
      </div>
    );
  }

  return (
    <div className="houses-list">
      <h2>Your Houses</h2>
      <div className="list-container">
        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Address</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {houses.map((house) => (
              <tr key={house.id}>
                <td>{house.name}</td>
                <td>
                  {house.street} {house.number}, {house.zipCode} {house.city}, {house.country}
                </td>
                <td className="actions">
                  <button 
                    className="button edit" 
                    onClick={() => onEdit(house)}
                    title="Edit house"
                  >
                    Edit
                  </button>
                  <button 
                    className="button delete" 
                    onClick={() => onDelete(house.id)}
                    title="Delete house"
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

export default HousesList;