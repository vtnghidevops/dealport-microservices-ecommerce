import React from 'react';

export function ReviewCard({ name, review, avatar, selected = false }) {
  return (
    <div className={`card p-4 ${selected ? 'border-2 border-ocean-green' : ''}`}>
      <div className="flex items-start">
        <img src={avatar} alt={name} className="w-10 h-10 rounded-full" />
        <div className="ml-3">
          <div className="flex items-center">
            <h4 className="font-bold">{name}</h4>
            <div className="text-pending ml-2">★★★★★</div>
          </div>
          <p className="text-sm mt-2">"{review}"</p>
        </div>
      </div>
      <div className="text-right mt-2 text-xs">
        {selected ? 'Selected' : 'Normal'}
      </div>
    </div>
  );
}