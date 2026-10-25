"use client";

import React, { useState, useEffect } from "react";
import Link from "next/link";

interface PatientRecord {
  id: string;
  displayName: string;
  identityType?: "NIC" | "PASSPORT";
  identityValue?: string;
  phone?: string;
  age: number;
  gender: "Male" | "Female" | "Other";
  city: string;
  dosha: "Vata-Pitta" | "Pitta-Kapha" | "Kapha-Vata" | "Tridoshic" | "Vata" | "Pitta" | "Kapha";
  concern: string;
  registeredAt: string;
  status: "Active" | "Under Treatment" | "Discharged";
}

const initialPatients: PatientRecord[] = [
  {
    id: "AYUR-2026-0142",
    displayName: "K. M. Bandara (කේ. එම්. බණ්ඩාර)",
    identityType: "NIC",
    identityValue: "198814200450",
    phone: "077-8912345",
    age: 46,
    gender: "Male",
    city: "Colombo",
    dosha: "Vata-Pitta",
    concern: "Sandhigata Vata (Osteoarthritis of the knees)",
    registeredAt: "2026-10-18 · 09:15 AM",
    status: "Under Treatment",
  },
  {
    id: "AYUR-2026-0128",
    displayName: "S. D. Jayawardena",
    identityType: "NIC",
    identityValue: "199468200192",
    phone: "071-2345678",
    age: 32,
    gender: "Female",
    city: "Nugegoda",
    dosha: "Pitta-Kapha",
    concern: "Amlapitta & Digestive Weakness (Acid Peptic Disease)",
    registeredAt: "2026-10-14 · 10:30 AM",
    status: "Under Treatment",
  },
  {
    id: "AYUR-2026-0105",
    displayName: "Nimali Perera",
    identityType: "NIC",
    identityValue: "197951200843",
    phone: "076-3456789",
    age: 55,
    gender: "Female",
    city: "Dehiwala",
    dosha: "Kapha-Vata",
    concern: "Kati Graha (Chronic Lumbar Stiffness)",
    registeredAt: "2026-10-09 · 08:45 AM",
    status: "Under Treatment",
  },
  {
    id: "AYUR-2026-0099",
    displayName: "Chandana Silva (චන්දන සිල්වා)",
    identityType: "NIC",
    identityValue: "198219400231",
    phone: "072-4567890",
    age: 44,
    gender: "Male",
    city: "Moratuwa",
    dosha: "Tridoshic",
    concern: "General Rejuvenation & Rasayana Therapy",
    registeredAt: "2026-10-04 · 11:10 AM",
    status: "Active",
  },
  {
    id: "AYUR-2026-0087",
    displayName: "Anoma Weerasinghe",
    identityType: "NIC",
    identityValue: "196851200991",
    phone: "075-5678901",
    age: 58,
    gender: "Female",
    city: "Kaduwela",
    dosha: "Vata-Pitta",
    concern: "Janu Sandhivata (Bilateral Joint Pain)",
    registeredAt: "2026-09-28 · 09:00 AM",
    status: "Active",
  },
];

export default function PatientsPage() {
  const [activeTab, setActiveTab] = useState<"register" | "directory">("register");
  const [patients, setPatients] = useState<PatientRecord[]>(initialPatients);
  const [currentTime, setCurrentTime] = useState("");

  // Live Asia/Colombo clock matching dashboard
  useEffect(() => {
    const updateTime = () => {
      const now = new Date();
      const timeStr = now.toLocaleTimeString("en-LK", {
        timeZone: "Asia/Colombo",
        hour: "2-digit",
        minute: "2-digit",
        second: "2-digit",
        hour12: true,
      });
      const dateStr = now.toLocaleDateString("en-LK", {
        timeZone: "Asia/Colombo",
        weekday: "short",
        day: "numeric",
        month: "short",
        year: "numeric",
      });
      setCurrentTime(`${dateStr} · ${timeStr} (SLST)`);
    };

    updateTime();
    const timer = setInterval(updateTime, 1000);
    return () => clearInterval(timer);
  }, []);

  // Form states
  const [displayName, setDisplayName] = useState("");
  const [hasIdentity, setHasIdentity] = useState(true);
  const [identityType, setIdentityType] = useState<"NIC" | "PASSPORT">("NIC");
  const [identityValue, setIdentityValue] = useState("");
  const [phone, setPhone] = useState("");
  const [age, setAge] = useState<string>("");
  const [gender, setGender] = useState<"Male" | "Female" | "Other">("Male");
  const [city, setCity] = useState("Colombo");
  const [dosha, setDosha] = useState<PatientRecord["dosha"]>("Vata-Pitta");
  const [concern, setConcern] = useState("");
  const [emergencyContact, setEmergencyContact] = useState("");
  const [emergencyPhone, setEmergencyPhone] = useState("");

  // Validation & UI states
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [successRecord, setSuccessRecord] = useState<PatientRecord | null>(null);
  const [directorySearch, setDirectorySearch] = useState("");

  const validate = () => {
    const errs: Record<string, string> = {};

    if (!displayName.trim()) {
      errs.displayName = "Full patient display name is required.";
    } else if (displayName.trim().length > 512) {
      errs.displayName = "Name must not exceed 512 characters.";
    }

    if (hasIdentity) {
      if (!identityValue.trim()) {
        errs.identityValue = `Please enter the ${identityType} number or uncheck identity if not provided.`;
      } else {
        const canonical = identityValue.replace(/[\s-]/g, "").toUpperCase();
        if (identityType === "NIC") {
          const isNewNic = /^\d{12}$/.test(canonical);
          const isOldNic = /^\d{9}[VX]$/.test(canonical);
          if (!isNewNic && !isOldNic) {
            errs.identityValue = "Invalid Sri Lankan NIC format (e.g. 198814200450 or 881420450V).";
          }
        } else if (identityType === "PASSPORT") {
          if (canonical.length < 5 || canonical.length > 20) {
            errs.identityValue = "Passport number must be between 5 and 20 characters.";
          }
        }

        const isDuplicate = patients.some(
          (p) =>
            p.identityValue &&
            p.identityValue.replace(/[\s-]/g, "").toUpperCase() === canonical
        );
        if (isDuplicate) {
          errs.identityValue = "A patient with this canonical identity is already registered (409 Conflict).";
        }
      }
    }

    if (phone.trim()) {
      const cleanPhone = phone.replace(/[\s-]/g, "");
      if (!/^(\+94|0)\d{9}$/.test(cleanPhone)) {
        errs.phone = "Enter a valid Sri Lankan phone number (e.g. 0771234567 or +94771234567).";
      }
    }

    if (age && (parseInt(age, 10) < 1 || parseInt(age, 10) > 125)) {
      errs.age = "Please enter a valid age between 1 and 125.";
    }

    setErrors(errs);
    return Object.keys(errs).length === 0;
  };

  const handleRegister = (e: React.FormEvent) => {
    e.preventDefault();
    if (!validate()) return;

    setIsSubmitting(true);

    setTimeout(() => {
      const newIndex = patients.length + 143;
      const newId = `AYUR-2026-0${newIndex}`;

      const now = new Date();
      const dateStr = now.toLocaleDateString("en-LK", {
        timeZone: "Asia/Colombo",
        year: "numeric",
        month: "2-digit",
        day: "2-digit",
      });
      const timeStr = now.toLocaleTimeString("en-LK", {
        timeZone: "Asia/Colombo",
        hour: "2-digit",
        minute: "2-digit",
        hour12: true,
      });

      const newPatient: PatientRecord = {
        id: newId,
        displayName: displayName.trim(),
        identityType: hasIdentity ? identityType : undefined,
        identityValue: hasIdentity ? identityValue.trim() : undefined,
        phone: phone.trim() || undefined,
        age: age ? parseInt(age, 10) : 35,
        gender,
        city: city.trim() || "Colombo",
        dosha,
        concern: concern.trim() || "Ayurvedic Assessment & Wellness Consultation",
        registeredAt: `${dateStr} · ${timeStr}`,
        status: "Active",
      };

      setPatients([newPatient, ...patients]);
      setSuccessRecord(newPatient);
      setIsSubmitting(false);

      setDisplayName("");
      setIdentityValue("");
      setPhone("");
      setAge("");
      setConcern("");
      setEmergencyContact("");
      setEmergencyPhone("");
      setErrors({});
    }, 450);
  };

  const filteredDirectory = patients.filter((p) => {
    if (!directorySearch.trim()) return true;
    const q = directorySearch.toLowerCase().trim();
    const matchesName = p.displayName.toLowerCase().includes(q);
    const matchesId = p.id.toLowerCase().includes(q);
    const matchesNic = p.identityValue ? p.identityValue.toLowerCase().includes(q) : false;
    const matchesPhone = p.phone ? p.phone.includes(q) : false;
    const matchesCity = p.city.toLowerCase().includes(q);
    return matchesName || matchesId || matchesNic || matchesPhone || matchesCity;
  });

  return (
    <div className="space-y-6 sm:space-y-8 font-sans">
      {/* ── Sub-Navigation / Header ── */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-green-900/10 pb-4">
        <div>
          <span className="text-xs font-bold uppercase tracking-wider text-emerald-800">
            Clinic Staff Workspace
          </span>
          <h1 className="text-2xl sm:text-3xl font-bold tracking-tight text-[#183a2d]">
            Patient Registration & Intake
          </h1>
        </div>

        {/* Local time badge */}
        <div className="flex items-center gap-2 self-start md:self-auto bg-white border border-green-900/15 px-3 py-1.5 rounded-xl shadow-xs text-xs text-slate-700 font-mono">
          <span className="w-2 h-2 rounded-full bg-emerald-600 animate-pulse" />
          <span>{currentTime || "Asia/Colombo · UTC+05:30"}</span>
        </div>
      </div>


      {/* ── Top Hero Greeting Banner (Botanical Theme like Dashboard) ── */}
      <div className="relative overflow-hidden rounded-3xl bg-gradient-to-br from-[#0c1f16] via-[#122e21] to-[#183a2d] text-white p-6 sm:p-8 shadow-xl border border-emerald-900/40">
        <div className="absolute -right-16 -top-16 w-80 h-80 rounded-full bg-emerald-500/10 blur-3xl pointer-events-none" />
        <div className="absolute right-40 -bottom-20 w-64 h-64 rounded-full bg-teal-400/10 blur-2xl pointer-events-none" />

        <div className="relative z-10 flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6">
          <div className="space-y-2">
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-emerald-500/15 border border-emerald-400/25 text-emerald-300 text-xs font-semibold uppercase tracking-wider">
              <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
              Patient Registration & Intake · AT02
            </div>
            <h2 className="text-2xl sm:text-3xl font-extrabold tracking-tight text-white">
              Ayurvedic Patient Registration
            </h2>
            <p className="text-emerald-100/75 text-sm sm:text-base max-w-2xl leading-relaxed">
              Register new clinic patients with full Unicode support (Sinhala, Tamil, English), optional canonical NIC/passport uniqueness, and preliminary Dosha intake profile.
            </p>
          </div>

          {/* View Switcher Pills */}
          <div className="flex items-center gap-2 bg-white/10 p-1.5 rounded-2xl border border-white/15 backdrop-blur-xs text-xs font-semibold shrink-0">
            <button
              type="button"
              onClick={() => {
                setActiveTab("register");
                setSuccessRecord(null);
              }}
              className={`px-4 py-2 rounded-xl transition-all ${
                activeTab === "register"
                  ? "bg-white text-emerald-950 font-bold shadow-xs"
                  : "text-emerald-100 hover:text-white hover:bg-white/10"
              }`}
            >
              📝 New Registration
            </button>
            <button
              type="button"
              onClick={() => setActiveTab("directory")}
              className={`px-4 py-2 rounded-xl transition-all flex items-center gap-1.5 ${
                activeTab === "directory"
                  ? "bg-white text-emerald-950 font-bold shadow-xs"
                  : "text-emerald-100 hover:text-white hover:bg-white/10"
              }`}
            >
              <span>🔍 Patient Directory</span>
              <span
                className={`text-[10px] px-2 py-0.5 rounded-full font-bold ${
                  activeTab === "directory" ? "bg-emerald-900 text-emerald-100" : "bg-emerald-950/80 text-emerald-300"
                }`}
              >
                {patients.length}
              </span>
            </button>
          </div>
        </div>

        <div className="mt-6 pt-5 border-t border-emerald-800/40 flex flex-wrap items-center justify-between gap-4 text-xs text-emerald-200/80">
          <div className="flex items-center gap-2">
            <span className="text-base">🌿</span>
            <span>
              <strong className="text-emerald-100">DEC02 Protocol:</strong> Identity uniqueness is enforced to prevent duplicate race collisions. Stored in UTC, displayed in Asia/Colombo.
            </span>
          </div>
          <div className="flex items-center gap-4 text-emerald-300 font-mono text-[11px]">
            <span>Total Patients: {patients.length}</span>
            <span>•</span>
            <span>Single Clinic Registry</span>
          </div>
        </div>
      </div>

      {/* ── Success Banner upon Registration ── */}
      {successRecord && activeTab === "register" && (
        <div className="rounded-3xl bg-gradient-to-r from-emerald-900 via-emerald-800 to-teal-900 text-white p-6 shadow-xl border border-emerald-500/30 animate-in fade-in zoom-in-95 duration-200">
          <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
            <div className="flex items-start gap-4">
              <div className="w-12 h-12 rounded-2xl bg-emerald-400/20 border border-emerald-400/40 flex items-center justify-center text-emerald-300 text-xl font-bold shrink-0">
                ✓
              </div>
              <div>
                <span className="text-[11px] font-bold uppercase tracking-widest text-emerald-300 bg-emerald-950/60 px-2.5 py-0.5 rounded-full border border-emerald-400/30">
                  Registration Successful · AT02 Verified
                </span>
                <h3 className="text-xl font-bold text-white mt-1">
                  {successRecord.displayName}
                </h3>
                <div className="flex items-center gap-3 text-xs text-emerald-100/80 mt-1 flex-wrap font-mono">
                  <span>ID: <strong className="text-emerald-200">{successRecord.id}</strong></span>
                  <span>•</span>
                  <span>{successRecord.identityType || "No ID"}: {successRecord.identityValue || "None"}</span>
                  <span>•</span>
                  <span>Registered: {successRecord.registeredAt}</span>
                </div>
              </div>
            </div>

            <div className="flex items-center gap-2.5 self-end sm:self-center shrink-0">
              <button
                type="button"
                onClick={() => setSuccessRecord(null)}
                className="px-3 py-2 rounded-xl bg-white/10 hover:bg-white/20 text-white text-xs font-semibold transition-colors"
              >
                Register Another
              </button>
              <button
                type="button"
                onClick={() => setActiveTab("directory")}
                className="px-4 py-2 rounded-xl bg-emerald-400 hover:bg-emerald-300 text-emerald-950 text-xs font-bold shadow-md shadow-emerald-950/40 transition-colors"
              >
                View in Directory →
              </button>
            </div>
          </div>
        </div>
      )}

      {/* ── TAB 1: REGISTRATION FORM ── */}
      {activeTab === "register" && (
        <form onSubmit={handleRegister} noValidate className="space-y-6 sm:space-y-8">
          <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 sm:gap-8">
            {/* Left Column (7 cols): Primary Identity & Contact Details */}
            <div className="lg:col-span-7 space-y-6">
              <div className="rounded-3xl bg-white p-6 sm:p-7 border border-green-900/10 shadow-xs space-y-5">
                <div className="border-b border-slate-100 pb-3">
                  <h3 className="text-lg font-bold text-[#183a2d] flex items-center gap-2">
                    <span>1. Patient Demographics & Identity</span>
                    <span className="text-[10px] text-emerald-800 bg-emerald-50 px-2 py-0.5 rounded font-mono">
                      TASK-PAT-01
                    </span>
                  </h3>
                  <p className="text-xs text-slate-500 mt-0.5">
                    Full name supports Unicode (Sinhala, Tamil, English). Identity document is optional per DEC02.
                  </p>
                </div>

                {/* Display Name */}
                <div className="space-y-1.5">
                  <label htmlFor="displayName" className="block text-xs font-bold text-slate-700">
                    Full Patient Name / Display Name <span className="text-rose-500">*</span>
                  </label>
                  <input
                    id="displayName"
                    type="text"
                    value={displayName}
                    onChange={(e) => setDisplayName(e.target.value)}
                    placeholder="e.g. K. M. Bandara or කේ. එම්. බණ්ඩාර"
                    className={`w-full text-sm px-3.5 py-2.5 rounded-xl border ${
                      errors.displayName
                        ? "border-rose-400 bg-rose-50/40 focus:outline-rose-500"
                        : "border-slate-200 bg-[#fbfcf9] focus:outline-emerald-600 focus:bg-white"
                    } transition-colors`}
                    required
                  />
                  {errors.displayName ? (
                    <p className="text-xs text-rose-600 font-medium">{errors.displayName}</p>
                  ) : (
                    <p className="text-[11px] text-slate-400">
                      Supports full Unicode text preserving Sinhala/Tamil script without truncation.
                    </p>
                  )}
                </div>

                {/* Identity Checkbox & Inputs */}
                <div className="p-4 rounded-2xl bg-[#f8faf7] border border-slate-200/80 space-y-3">
                  <div className="flex items-center justify-between">
                    <label className="flex items-center gap-2 text-xs font-bold text-slate-700 cursor-pointer select-none">
                      <input
                        type="checkbox"
                        checked={hasIdentity}
                        onChange={(e) => setHasIdentity(e.target.checked)}
                        className="w-4 h-4 rounded text-emerald-700 accent-emerald-700"
                      />
                      <span>Patient has Identity Document (NIC / Passport)</span>
                    </label>
                    <span className="text-[10px] font-semibold text-slate-400 uppercase tracking-wider">
                      DEC02 Optional
                    </span>
                  </div>

                  {hasIdentity && (
                    <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 pt-1">
                      <div>
                        <label htmlFor="identityType" className="block text-[11px] font-bold text-slate-600 mb-1">
                          Document Type
                        </label>
                        <select
                          id="identityType"
                          value={identityType}
                          onChange={(e) => setIdentityType(e.target.value as "NIC" | "PASSPORT")}
                          className="w-full text-xs px-3 py-2 rounded-xl border border-slate-200 bg-white text-slate-800 focus:outline-emerald-600 font-medium"
                        >
                          <option value="NIC">National Identity Card (NIC)</option>
                          <option value="PASSPORT">Passport</option>
                        </select>
                      </div>

                      <div className="sm:col-span-2">
                        <label htmlFor="identityValue" className="block text-[11px] font-bold text-slate-600 mb-1">
                          Document Number <span className="text-rose-500">*</span>
                        </label>
                        <input
                          id="identityValue"
                          type="text"
                          value={identityValue}
                          onChange={(e) => setIdentityValue(e.target.value)}
                          placeholder={identityType === "NIC" ? "198814200450 or 881420450V" : "e.g. N1234567"}
                          className={`w-full text-xs px-3 py-2 rounded-xl border ${
                            errors.identityValue
                              ? "border-rose-400 bg-rose-50/40 focus:outline-rose-500"
                              : "border-slate-200 bg-white focus:outline-emerald-600"
                          } font-mono`}
                        />
                      </div>
                    </div>
                  )}

                  {hasIdentity && errors.identityValue && (
                    <p className="text-xs text-rose-600 font-medium">{errors.identityValue}</p>
                  )}

                  <p className="text-[11px] text-slate-500 leading-relaxed">
                    Identities are normalized into canonical alphanumeric uppercase in the database to prevent duplicate registration collisions.
                  </p>
                </div>

                {/* Contact Phone & City */}
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div className="space-y-1">
                    <label htmlFor="phone" className="block text-xs font-bold text-slate-700">
                      Mobile / Contact Phone
                    </label>
                    <input
                      id="phone"
                      type="tel"
                      value={phone}
                      onChange={(e) => setPhone(e.target.value)}
                      placeholder="e.g. 077 123 4567"
                      className={`w-full text-xs px-3.5 py-2.5 rounded-xl border ${
                        errors.phone
                          ? "border-rose-400 bg-rose-50/40 focus:outline-rose-500"
                          : "border-slate-200 bg-[#fbfcf9] focus:outline-emerald-600 focus:bg-white"
                      } font-mono`}
                    />
                    {errors.phone && <p className="text-xs text-rose-600">{errors.phone}</p>}
                  </div>

                  <div className="space-y-1">
                    <label htmlFor="city" className="block text-xs font-bold text-slate-700">
                      City / District
                    </label>
                    <input
                      id="city"
                      type="text"
                      value={city}
                      onChange={(e) => setCity(e.target.value)}
                      placeholder="e.g. Colombo, Kandy, Galle"
                      className="w-full text-xs px-3.5 py-2.5 rounded-xl border border-slate-200 bg-[#fbfcf9] focus:outline-emerald-600 focus:bg-white"
                    />
                  </div>
                </div>

                {/* Age & Gender */}
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div className="space-y-1">
                    <label htmlFor="age" className="block text-xs font-bold text-slate-700">
                      Age (Years)
                    </label>
                    <input
                      id="age"
                      type="number"
                      min="1"
                      max="125"
                      value={age}
                      onChange={(e) => setAge(e.target.value)}
                      placeholder="e.g. 42"
                      className={`w-full text-xs px-3.5 py-2.5 rounded-xl border ${
                        errors.age
                          ? "border-rose-400 bg-rose-50/40 focus:outline-rose-500"
                          : "border-slate-200 bg-[#fbfcf9] focus:outline-emerald-600 focus:bg-white"
                      }`}
                    />
                    {errors.age && <p className="text-xs text-rose-600">{errors.age}</p>}
                  </div>

                  <div className="space-y-1">
                    <label className="block text-xs font-bold text-slate-700">Gender</label>
                    <div className="flex items-center gap-2 pt-1">
                      {(["Male", "Female", "Other"] as const).map((g) => (
                        <label
                          key={g}
                          className={`flex-1 flex items-center justify-center py-2 px-3 rounded-xl border text-xs font-semibold cursor-pointer transition-all ${
                            gender === g
                              ? "bg-emerald-50 border-emerald-600 text-emerald-900 shadow-xs"
                              : "border-slate-200 bg-[#fbfcf9] text-slate-600 hover:bg-slate-100"
                          }`}
                        >
                          <input
                            type="radio"
                            name="gender"
                            value={g}
                            checked={gender === g}
                            onChange={() => setGender(g)}
                            className="sr-only"
                          />
                          <span>{g}</span>
                        </label>
                      ))}
                    </div>
                  </div>
                </div>
              </div>

              {/* Emergency Contact */}
              <div className="rounded-3xl bg-white p-6 border border-green-900/10 shadow-xs space-y-4">
                <h3 className="text-sm font-bold text-[#183a2d]">Emergency Contact Person (Optional)</h3>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div>
                    <label htmlFor="emergencyContact" className="block text-xs font-semibold text-slate-600 mb-1">
                      Contact Name
                    </label>
                    <input
                      id="emergencyContact"
                      type="text"
                      value={emergencyContact}
                      onChange={(e) => setEmergencyContact(e.target.value)}
                      placeholder="e.g. S. M. Bandara (Spouse)"
                      className="w-full text-xs px-3.5 py-2.5 rounded-xl border border-slate-200 bg-[#fbfcf9] focus:outline-emerald-600 focus:bg-white"
                    />
                  </div>
                  <div>
                    <label htmlFor="emergencyPhone" className="block text-xs font-semibold text-slate-600 mb-1">
                      Emergency Phone
                    </label>
                    <input
                      id="emergencyPhone"
                      type="tel"
                      value={emergencyPhone}
                      onChange={(e) => setEmergencyPhone(e.target.value)}
                      placeholder="e.g. 071 987 6543"
                      className="w-full text-xs px-3.5 py-2.5 rounded-xl border border-slate-200 bg-[#fbfcf9] focus:outline-emerald-600 focus:bg-white font-mono"
                    />
                  </div>
                </div>
              </div>
            </div>

            {/* Right Column (5 cols): Ayurvedic Intake & Submit */}
            <div className="lg:col-span-5 space-y-6">
              <div className="rounded-3xl bg-white p-6 sm:p-7 border border-green-900/10 shadow-xs space-y-5">
                <div className="border-b border-slate-100 pb-3">
                  <h3 className="text-lg font-bold text-[#183a2d] flex items-center gap-2">
                    <span>2. Ayurvedic Intake Assessment</span>
                    <span className="text-[10px] text-teal-800 bg-teal-50 px-2 py-0.5 rounded font-mono">
                      DEC02
                    </span>
                  </h3>
                  <p className="text-xs text-slate-500 mt-0.5">
                    Provisional dosha prakriti and initial presenting concerns for clinician review.
                  </p>
                </div>

                {/* Dosha Selector */}
                <div className="space-y-2">
                  <label className="block text-xs font-bold text-slate-700">
                    Provisional Prakriti / Dosha Profile
                  </label>
                  <div className="grid grid-cols-2 gap-2">
                    {[
                      { val: "Vata-Pitta", desc: "Variable digestion & mobility" },
                      { val: "Pitta-Kapha", desc: "Moderate metabolic strength" },
                      { val: "Kapha-Vata", desc: "Cooling constitution" },
                      { val: "Tridoshic", desc: "Balanced three humours" },
                      { val: "Vata", desc: "Primary Vata dominance" },
                      { val: "Pitta", desc: "Primary Pitta dominance" },
                      { val: "Kapha", desc: "Primary Kapha dominance" },
                    ].map((d) => (
                      <button
                        type="button"
                        key={d.val}
                        onClick={() => setDosha(d.val as PatientRecord["dosha"])}
                        className={`p-2.5 rounded-xl border text-left transition-all ${
                          dosha === d.val
                            ? "bg-gradient-to-r from-emerald-50 to-teal-50 border-emerald-600 text-emerald-950 shadow-xs"
                            : "border-slate-200 bg-[#fbfcf9] hover:bg-slate-50 text-slate-700"
                        }`}
                      >
                        <p className="text-xs font-bold">{d.val}</p>
                        <p className="text-[10px] text-slate-500 truncate">{d.desc}</p>
                      </button>
                    ))}
                  </div>
                </div>

                {/* Chief Health Concern */}
                <div className="space-y-1.5">
                  <label htmlFor="concern" className="block text-xs font-bold text-slate-700">
                    Primary Health Concern / Presenting Symptoms
                  </label>
                  <textarea
                    id="concern"
                    rows={3}
                    value={concern}
                    onChange={(e) => setConcern(e.target.value)}
                    placeholder="e.g. Lumbar stiffness (Kati Graha) aggravated in cold weather, irregular digestion (Vishamagni)..."
                    className="w-full text-xs p-3 rounded-xl border border-slate-200 bg-[#fbfcf9] focus:outline-emerald-600 focus:bg-white leading-relaxed resize-none"
                  />
                  <p className="text-[11px] text-slate-400">
                    Consultation notes remain draft until clinician signs (AT03 immutability).
                  </p>
                </div>

                {/* Audit & Data Governance Notice */}
                <div className="p-4 rounded-2xl bg-[#0c1f16] text-emerald-200 text-xs border border-emerald-900/60 space-y-1.5">
                  <div className="flex items-center gap-2 text-white font-bold text-xs">
                    <span className="w-2 h-2 rounded-full bg-emerald-400" />
                    <span>Transactional Audit & Integrity (AT13)</span>
                  </div>
                  <p className="text-[11px] text-emerald-300/80 leading-relaxed">
                    Patient creation inserts an audit event atomically in PostgreSQL. Clinical details are restricted to authorized clinical staff.
                  </p>
                </div>

                {/* Submit Action Buttons */}
                <div className="pt-2 space-y-2">
                  <button
                    type="submit"
                    disabled={isSubmitting}
                    className="w-full py-3.5 px-4 rounded-xl bg-gradient-to-r from-[#23745b] to-[#10b981] hover:from-[#1e4d3a] hover:to-[#059669] text-white text-sm font-bold shadow-md shadow-emerald-900/30 transition-all flex items-center justify-center gap-2 disabled:opacity-50 cursor-pointer"
                  >
                    {isSubmitting ? (
                      <>
                        <span className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                        <span>Registering Patient...</span>
                      </>
                    ) : (
                      <>
                        <svg className="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                          <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
                          <circle cx="9" cy="7" r="4" />
                          <line x1="19" x2="19" y1="8" y2="14" />
                          <line x1="22" x2="16" y1="11" y2="11" />
                        </svg>
                        <span>Register Patient (Generate AYUR ID)</span>
                      </>
                    )}
                  </button>

                  <button
                    type="button"
                    onClick={() => {
                      setDisplayName("");
                      setIdentityValue("");
                      setPhone("");
                      setAge("");
                      setConcern("");
                      setErrors({});
                    }}
                    className="w-full py-2 text-xs font-semibold text-slate-500 hover:text-slate-800 transition-colors cursor-pointer"
                  >
                    Clear Form
                  </button>
                </div>
              </div>
            </div>
          </div>
        </form>
      )}

      {/* ── TAB 2: PATIENT DIRECTORY & LOOKUP ── */}
      {activeTab === "directory" && (
        <div className="rounded-3xl bg-white border border-green-900/10 shadow-xs overflow-hidden space-y-0">
          <div className="p-5 sm:p-6 border-b border-slate-100 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
            <div>
              <h3 className="text-lg sm:text-xl font-bold text-[#183a2d]">
                Verified Patient Directory (TASK-PAT-02 / AT02)
              </h3>
              <p className="text-xs text-slate-500 mt-0.5">
                Search verified records by Unicode name, canonical NIC/passport, registration ID, or phone number.
              </p>
            </div>

            <button
              type="button"
              onClick={() => setActiveTab("register")}
              className="px-4 py-2 rounded-xl bg-[#23745b] hover:bg-[#183a2d] text-white text-xs font-semibold shadow-xs transition-colors self-start sm:self-auto flex items-center gap-1.5 cursor-pointer"
            >
              <svg className="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M12 5v14M5 12h14" />
              </svg>
              <span>Add New Patient</span>
            </button>
          </div>

          {/* Search Bar */}
          <div className="px-5 sm:px-6 py-3 bg-[#fbfcf9] border-b border-slate-100 flex items-center justify-between gap-4">
            <div className="relative flex-1 max-w-md">
              <svg className="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <circle cx="11" cy="11" r="8" />
                <path d="m21 21-4.3-4.3" />
              </svg>
              <input
                type="text"
                value={directorySearch}
                onChange={(e) => setDirectorySearch(e.target.value)}
                placeholder="Search by Unicode name (සිංහල), NIC, ID (AYUR-2026-XXXX)..."
                className="w-full text-xs pl-9 pr-3 py-2 rounded-xl border border-slate-200 bg-white text-slate-800 focus:outline-emerald-600 shadow-xs"
              />
            </div>
            <span className="text-xs text-slate-500 font-medium">
              Showing {filteredDirectory.length} of {patients.length} records
            </span>
          </div>

          {/* Patient Records List */}
          <div className="divide-y divide-slate-100 overflow-x-auto">
            {filteredDirectory.length === 0 ? (
              <div className="p-12 text-center text-slate-400">
                <p className="text-sm font-semibold">No patient records match your search query.</p>
                <p className="text-xs text-slate-500 mt-1">
                  Try searching by full name, NIC without spaces, or clinic ID.
                </p>
              </div>
            ) : (
              filteredDirectory.map((patient) => (
                <div
                  key={patient.id}
                  className="p-4 sm:p-5 hover:bg-emerald-50/40 transition-colors flex flex-col lg:flex-row lg:items-center justify-between gap-4"
                >
                  <div className="flex items-start gap-3.5">
                    <div className="w-10 h-10 rounded-xl bg-emerald-100 text-emerald-900 border border-emerald-300 flex items-center justify-center font-bold text-xs shrink-0">
                      {patient.displayName
                        .split(" ")
                        .map((n) => n[0])
                        .slice(0, 2)
                        .join("")}
                    </div>

                    <div>
                      <div className="flex items-center gap-2 flex-wrap">
                        <span className="font-bold text-sm text-[#183a2d]">
                          {patient.displayName}
                        </span>
                        <span className="font-mono text-[11px] text-slate-600 bg-slate-100 px-2 py-0.5 rounded font-semibold">
                          {patient.id}
                        </span>
                        <span className="text-[10px] font-semibold bg-emerald-100 text-emerald-800 px-2 py-0.5 rounded-full">
                          {patient.dosha}
                        </span>
                        <span className="text-[10px] font-medium text-slate-500">
                          {patient.gender}, {patient.age} yrs · {patient.city}
                        </span>
                      </div>

                      <p className="text-xs text-slate-700 font-medium mt-1">
                        {patient.concern}
                      </p>

                      <div className="flex items-center gap-3 text-xs text-slate-500 mt-1 flex-wrap font-mono">
                        {patient.identityValue && (
                          <span>
                            {patient.identityType}: {patient.identityValue}
                          </span>
                        )}
                        {patient.phone && (
                          <>
                            <span>•</span>
                            <span>Phone: {patient.phone}</span>
                          </>
                        )}
                        <span>•</span>
                        <span>Registered: {patient.registeredAt}</span>
                      </div>
                    </div>
                  </div>

                  <div className="flex items-center gap-2 shrink-0 self-end lg:self-center">
                    <span
                      className={`text-[11px] font-bold px-2.5 py-1 rounded-full ${
                        patient.status === "Under Treatment"
                          ? "bg-amber-100 text-amber-900 border border-amber-300"
                          : "bg-emerald-100 text-emerald-900 border border-emerald-300"
                      }`}
                    >
                      {patient.status}
                    </span>

                    <Link
                      href="/clinical"
                      className="px-3 py-1.5 text-xs font-semibold rounded-lg bg-[#23745b] hover:bg-[#183a2d] text-white shadow-xs transition-colors"
                    >
                      Care Pathway →
                    </Link>
                  </div>
                </div>
              ))
            )}
          </div>

          <div className="p-4 bg-[#fbfcf9] border-t border-slate-100 flex items-center justify-between text-xs text-slate-500">
            <span>Identity uniqueness enforced at PostgreSQL repository layer (Migration 000005).</span>
            <span className="font-mono text-[11px]">Sprint 1 Verified Slices</span>
          </div>
        </div>
      )}
    </div>
  );
}
