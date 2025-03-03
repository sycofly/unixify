// JavaScript for the section pages (people, system, database, service)
document.addEventListener('DOMContentLoaded', function() {
    // Get section type from the hidden input
    const sectionType = document.getElementById('accountType').value;
    
    // Add section-specific class to body for styling
    document.body.classList.add(`section-${sectionType}`);
    
    // Elements
    const accountsTableBody = document.getElementById('accountsTableBody');
    const groupsTableBody = document.getElementById('groupsTableBody');
    const auditTableBody = document.getElementById('auditTableBody');
    const auditPagination = document.getElementById('auditPagination');
    const newAccountBtn = document.getElementById('newAccountBtn');
    const newGroupBtn = document.getElementById('newGroupBtn');
    const saveAccountBtn = document.getElementById('saveAccountBtn');
    const saveGroupBtn = document.getElementById('saveGroupBtn');
    const searchInput = document.getElementById('searchInput');
    const searchBtn = document.getElementById('searchBtn');
    const refreshAuditBtn = document.getElementById('refreshAuditBtn');
    const auditFilterEntity = document.getElementById('auditFilterEntity');
    const auditFilterAction = document.getElementById('auditFilterAction');
    
    // Audit state
    let auditCurrentPage = 1;
    let auditTotalPages = 1;
    let auditPageSize = 10;
    
    // Bootstrap modals
    let accountModal, groupModal, membershipModal, auditDetailModal;
    if (document.getElementById('accountModal')) {
        accountModal = new bootstrap.Modal(document.getElementById('accountModal'));
    }
    if (document.getElementById('groupModal')) {
        groupModal = new bootstrap.Modal(document.getElementById('groupModal'));
    }
    if (document.getElementById('membershipModal')) {
        membershipModal = new bootstrap.Modal(document.getElementById('membershipModal'));
    }
    if (document.getElementById('auditDetailModal')) {
        auditDetailModal = new bootstrap.Modal(document.getElementById('auditDetailModal'));
    }
    
    // Load accounts and groups on page load
    loadAccounts();
    loadGroups();
    
    // Load audit logs when the tab is clicked
    document.getElementById('audit-tab').addEventListener('shown.bs.tab', function() {
        loadAuditLogs();
    });
    
    // Event listeners
    if (newAccountBtn) {
        newAccountBtn.addEventListener('click', openNewAccountModal);
    }
    if (newGroupBtn) {
        newGroupBtn.addEventListener('click', openNewGroupModal);
    }
    if (saveAccountBtn) {
        saveAccountBtn.addEventListener('click', saveAccount);
    }
    if (saveGroupBtn) {
        saveGroupBtn.addEventListener('click', saveGroup);
    }
    if (searchBtn && searchInput) {
        searchBtn.addEventListener('click', () => {
            performSearch(searchInput.value);
        });
        searchInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                performSearch(searchInput.value);
            }
        });
    }
    if (refreshAuditBtn) {
        refreshAuditBtn.addEventListener('click', loadAuditLogs);
    }
    if (auditFilterEntity) {
        auditFilterEntity.addEventListener('change', loadAuditLogs);
    }
    if (auditFilterAction) {
        auditFilterAction.addEventListener('change', loadAuditLogs);
    }
    
    // Setup UID/GID validation
    setupValidation();
    
    // Functions
    
    // Set up input validation based on section type
    function setupValidation() {
        const uidField = document.getElementById('uid');
        const gidField = document.getElementById('gid');
        
        if (uidField) {
            uidField.addEventListener('change', validateUID);
        }
        
        if (gidField) {
            gidField.addEventListener('change', validateGID);
        }
    }
    
    // Check for duplicate UID
    async function checkUIDDuplicate(uid, accountId = null) {
        try {
            let url = `/api/accounts/check-duplicate?uid=${uid}`;
            if (accountId) {
                url += `&exclude_id=${accountId}`;
            }
            const response = await authFetch(url);
            const data = await response.json();
            return data.isDuplicate === true;
        } catch (error) {
            console.error('Error checking for duplicate UID:', error);
            return false; // Assume no duplicate if API call fails
        }
    }
    
    // Validate UID based on account type
    async function validateUID() {
        const uidField = document.getElementById('uid');
        const uid = parseInt(uidField.value);
        const accountId = document.getElementById('accountId').value || null;
        const uidFeedback = document.getElementById('uidFeedback') || 
            document.createElement('div');
        
        if (!document.getElementById('uidFeedback')) {
            uidFeedback.id = 'uidFeedback';
            uidFeedback.className = 'form-text';
            uidField.after(uidFeedback);
        }
        
        let isValid = true;
        let message = '';
        
        // Basic range validation (now just a warning)
        switch (sectionType) {
            case 'people':
                if (uid < 1000 || uid > 60000) {
                    message = 'Warning: People account UIDs are recommended to be between 1000 and 60000';
                }
                break;
            case 'system':
                if (uid < 9000 || uid > 9100) {
                    message = 'Warning: System account UIDs are recommended to be between 9000 and 9100';
                }
                break;
            case 'service':
                if (uid < 60001 || uid > 65535) {
                    message = 'Warning: Service account UIDs are recommended to be between 60001 and 65535';
                }
                break;
            case 'database':
                if (uid < 70000 || uid > 79999) {
                    message = 'Warning: Database account UIDs are recommended to be between 70000 and 79999';
                }
                break;
        }
        
        // Check for duplicate UID
        const isDuplicate = await checkUIDDuplicate(uid, accountId);
        if (isDuplicate) {
            isValid = false;
            message = `Error: UID ${uid} is already in use by another account`;
            uidField.style.backgroundColor = '#ffdddd';  // Light red background
        } else {
            uidField.style.backgroundColor = '';  // Reset background
        }
        
        if (!isValid) {
            uidField.classList.add('is-invalid');
            uidField.classList.remove('is-valid');
            uidFeedback.className = 'invalid-feedback d-block';
            uidFeedback.textContent = message;
        } else if (message) {
            // It's just a warning, not an error
            uidField.classList.add('is-valid');
            uidField.classList.remove('is-invalid');
            uidFeedback.className = 'text-warning d-block';
            uidFeedback.textContent = message;
        } else {
            uidField.classList.remove('is-invalid');
            uidField.classList.add('is-valid');
            uidFeedback.textContent = '';
        }
        
        return isValid;
    }
    
    // Check for duplicate GID
    async function checkGIDDuplicate(gid, groupId = null) {
        try {
            let url = `/api/groups/check-duplicate?gid=${gid}`;
            if (groupId) {
                url += `&exclude_id=${groupId}`;
            }
            const response = await authFetch(url);
            const data = await response.json();
            return data.isDuplicate === true;
        } catch (error) {
            console.error('Error checking for duplicate GID:', error);
            return false; // Assume no duplicate if API call fails
        }
    }
    
    // Validate GID based on group type
    async function validateGID() {
        const gidField = document.getElementById('gid');
        const gid = parseInt(gidField.value);
        const groupId = document.getElementById('groupId').value || null;
        const gidFeedback = document.getElementById('gidFeedback') || 
            document.createElement('div');
        
        if (!document.getElementById('gidFeedback')) {
            gidFeedback.id = 'gidFeedback';
            gidFeedback.className = 'form-text';
            gidField.after(gidFeedback);
        }
        
        let isValid = true;
        let message = '';
        
        // Basic range validation (now just a warning)
        switch (sectionType) {
            case 'people':
                if (gid < 1000 || gid > 60000) {
                    message = 'Warning: People group GIDs are recommended to be between 1000 and 60000';
                }
                break;
            case 'system':
                if (gid < 9000 || gid > 9100) {
                    message = 'Warning: System group GIDs are recommended to be between 9000 and 9100';
                }
                break;
            case 'service':
                if (gid < 60001 || gid > 65535) {
                    message = 'Warning: Service group GIDs are recommended to be between 60001 and 65535';
                }
                break;
            case 'database':
                if (gid < 70000 || gid > 79999) {
                    message = 'Warning: Database group GIDs are recommended to be between 70000 and 79999';
                }
                break;
        }
        
        // Check for duplicate GID
        const isDuplicate = await checkGIDDuplicate(gid, groupId);
        if (isDuplicate) {
            isValid = false;
            message = `Error: GID ${gid} is already in use by another group`;
            gidField.style.backgroundColor = '#ffdddd';  // Light red background
        } else {
            gidField.style.backgroundColor = '';  // Reset background
        }
        
        if (!isValid) {
            gidField.classList.add('is-invalid');
            gidField.classList.remove('is-valid');
            gidFeedback.className = 'invalid-feedback d-block';
            gidFeedback.textContent = message;
        } else if (message) {
            // It's just a warning, not an error
            gidField.classList.add('is-valid');
            gidField.classList.remove('is-invalid');
            gidFeedback.className = 'text-warning d-block';
            gidFeedback.textContent = message;
        } else {
            gidField.classList.remove('is-invalid');
            gidField.classList.add('is-valid');
            gidFeedback.textContent = '';
        }
        
        return isValid;
    }
    
    // Load accounts for the current section
    async function loadAccounts() {
        try {
            // Add a timestamp parameter to avoid caching issues
            const timestamp = new Date().getTime();
            // Use authFetch to include authorization token
            const response = await authFetch(`/api/accounts?type=${sectionType}&_=${timestamp}`);
            if (!response.ok) {
                throw new Error(`Error: ${response.status} ${response.statusText}`);
            }
            const accounts = await response.json();
            renderAccountsTable(accounts);
            return accounts;
        } catch (error) {
            console.error('Error loading accounts:', error);
            showError('Failed to load accounts: ' + error.message);
            return [];
        }
    }
    
    // Load groups for the current section
    async function loadGroups() {
        try {
            // Add a timestamp parameter to avoid caching issues
            const timestamp = new Date().getTime();
            // Use authFetch to include authorization token
            const response = await authFetch(`/api/groups?type=${sectionType}&_=${timestamp}`);
            if (!response.ok) {
                throw new Error(`Error: ${response.status} ${response.statusText}`);
            }
            const groups = await response.json();
            renderGroupsTable(groups);
            return groups;
        } catch (error) {
            console.error('Error loading groups:', error);
            showError('Failed to load groups: ' + error.message);
            return [];
        }
    }
    
    // Perform search across users and groups
    async function performSearch(query) {
        if (!query) {
            showError('Please enter a search term');
            return;
        }
        
        try {
            // Show loading indicators
            if (accountsTableBody) {
                accountsTableBody.innerHTML = '<tr><td colspan="6" class="text-center">Searching accounts...</td></tr>';
            }
            if (groupsTableBody) {
                groupsTableBody.innerHTML = '<tr><td colspan="6" class="text-center">Searching groups...</td></tr>';
            }
            
            // Search for accounts and groups simultaneously
            const [accounts, groups] = await Promise.all([
                apiRequest(`/api/search/accounts?q=${encodeURIComponent(query)}`),
                apiRequest(`/api/search/groups?q=${encodeURIComponent(query)}`)
            ]);
            
            // Filter results by current section type
            const filteredAccounts = accounts.filter(a => a.type.toLowerCase() === sectionType.toLowerCase());
            const filteredGroups = groups.filter(g => g.type.toLowerCase() === sectionType.toLowerCase());
            
            // Display results
            renderAccountsTable(filteredAccounts);
            renderGroupsTable(filteredGroups);
            
            // Show summary
            showSuccess(`Found ${filteredAccounts.length} accounts and ${filteredGroups.length} groups matching "${query}"`);
        } catch (error) {
            console.error('Error searching:', error);
            showError('Search failed: ' + error.message);
            
            // Reload original data
            loadAccounts();
            loadGroups();
        }
    }
    
    // Render accounts table
    function renderAccountsTable(accounts) {
        if (!accountsTableBody) return;
        
        accountsTableBody.innerHTML = '';
        
        if (!accounts || accounts.length === 0) {
            const row = document.createElement('tr');
            row.innerHTML = '<td colspan="6" class="text-center">No accounts found</td>';
            accountsTableBody.appendChild(row);
            return;
        }
        
        accounts.forEach(account => {
            // Get primary group info - with improved handling
            let primaryGroupText = 'None';
            if (account.primary_group && account.primary_group.groupname) {
                primaryGroupText = `${account.primary_group.groupname} (${account.primary_group.gid})`;
            } else if (account.primary_group_id) {
                // If we only have the ID but not the full group data, show the ID
                primaryGroupText = `ID: ${account.primary_group_id}`;
            }
            
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${account.id}</td>
                <td>${account.uid}</td>
                <td>${account.username}</td>
                <td>${primaryGroupText}</td>
                <td>${formatDate(account.created_at)}</td>
                <td class="action-buttons">
                    <button type="button" class="btn btn-sm btn-primary edit-account" data-id="${account.id}" title="Edit">
                        <i class="bi bi-pencil"></i> Edit
                    </button>
                    <button type="button" class="btn btn-sm btn-danger delete-account" data-id="${account.id}" title="Delete">
                        <i class="bi bi-trash"></i> Delete
                    </button>
                    <button type="button" class="btn btn-sm btn-info groups-account" data-id="${account.id}" title="Manage Groups">
                        <i class="bi bi-people"></i> Groups
                    </button>
                </td>
            `;
            accountsTableBody.appendChild(row);
        });
        
        // Add event listeners to buttons
        document.querySelectorAll('.edit-account').forEach(btn => {
            btn.addEventListener('click', () => editAccount(btn.dataset.id));
        });
        
        document.querySelectorAll('.delete-account').forEach(btn => {
            btn.addEventListener('click', () => deleteAccount(btn.dataset.id));
        });
        
        document.querySelectorAll('.groups-account').forEach(btn => {
            btn.addEventListener('click', () => manageAccountGroups(btn.dataset.id));
        });
    }
    
    // Render groups table
    function renderGroupsTable(groups) {
        if (!groupsTableBody) return;
        
        groupsTableBody.innerHTML = '';
        
        if (!groups || groups.length === 0) {
            const row = document.createElement('tr');
            row.innerHTML = '<td colspan="7" class="text-center">No groups found</td>';
            groupsTableBody.appendChild(row);
            return;
        }
        
        groups.forEach(group => {
            const row = document.createElement('tr');
            // Ensure description is displayed properly
            const displayDescription = group.description !== null && group.description !== undefined ? group.description : '';
            
            // Ensure created_by is displayed properly
            const displayCreatedBy = group.created_by !== null && group.created_by !== undefined ? group.created_by : '—';
            
            row.innerHTML = `
                <td>${group.id}</td>
                <td>${group.gid}</td>
                <td>${group.groupname}</td>
                <td>${displayDescription}</td>
                <td>${displayCreatedBy}</td>
                <td>${formatDate(group.created_at)}</td>
                <td class="action-buttons">
                    <button type="button" class="btn btn-sm btn-primary edit-group" data-id="${group.id}" title="Edit">
                        <i class="bi bi-pencil"></i> Edit
                    </button>
                    <button type="button" class="btn btn-sm btn-danger delete-group" data-id="${group.id}" title="Delete">
                        <i class="bi bi-trash"></i> Delete
                    </button>
                    <button type="button" class="btn btn-sm btn-info members-group" data-id="${group.id}" title="Manage Members">
                        <i class="bi bi-people"></i> Members
                    </button>
                </td>
            `;
            groupsTableBody.appendChild(row);
        });
        
        // Add event listeners to buttons
        document.querySelectorAll('.edit-group').forEach(btn => {
            btn.addEventListener('click', () => editGroup(btn.dataset.id));
        });
        
        document.querySelectorAll('.delete-group').forEach(btn => {
            btn.addEventListener('click', () => deleteGroup(btn.dataset.id));
        });
        
        document.querySelectorAll('.members-group').forEach(btn => {
            btn.addEventListener('click', () => manageGroupMembers(btn.dataset.id));
        });
    }
    
    // Open modal for new account
    async function openNewAccountModal() {
        if (!accountModal) return;
        
        document.getElementById('accountModalLabel').textContent = 'New Account';
        document.getElementById('accountForm').reset();
        document.getElementById('accountId').value = '';
        
        // Clear any previous validation indicators
        const uidField = document.getElementById('uid');
        if (uidField) {
            uidField.classList.remove('is-valid', 'is-invalid');
            uidField.style.backgroundColor = '';
            const feedback = document.getElementById('uidFeedback');
            if (feedback) feedback.textContent = '';
        }
        
        // Load groups for the primary group dropdown
        loadGroupsForDropdown();
        
        // Set up form requirements based on section type
        const currentSectionType = document.getElementById('accountType').value;
        const primaryGroupField = document.getElementById('primaryGroupID');
        const primaryGroupLabel = primaryGroupField.closest('.mb-3').querySelector('.form-label');
        
        // Set form attributes based on section type
        if (uidField) {
            switch (currentSectionType) {
                case 'people':
                    uidField.setAttribute('min', '1000');
                    uidField.setAttribute('max', '60000');
                    uidField.setAttribute('placeholder', '1000-60000 (recommended)');
                    break;
                case 'system':
                    uidField.setAttribute('min', '9000');
                    uidField.setAttribute('max', '9100');
                    uidField.setAttribute('placeholder', '9000-9100 (recommended)');
                    break;
                case 'service':
                    uidField.setAttribute('min', '60001');
                    uidField.setAttribute('max', '65535');
                    uidField.setAttribute('placeholder', '60001-65535 (recommended)');
                    break;
                case 'database':
                    uidField.setAttribute('min', '70000');
                    uidField.setAttribute('max', '79999');
                    uidField.setAttribute('placeholder', '70000-79999 (recommended)');
                    break;
            }
        }
        
        // Fetch the next available UID and set it as default value
        try {
            const response = await authFetch(`/api/accounts/next-uid?type=${currentSectionType}`);
            const data = await response.json();
            if (response.ok && data && data.uid) {
                uidField.value = data.uid;
                uidField.setAttribute('placeholder', `Suggestion: ${data.uid}`);
            }
        } catch (error) {
            console.warn('Error fetching next available UID:', error);
        }
        
        if (currentSectionType === 'system') {
            // For system accounts, primary group is required
            primaryGroupField.setAttribute('required', 'required');
            // Add a star to indicate required field
            primaryGroupLabel.innerHTML = 'Primary Group <span class="text-danger">*</span>';
            // Add helper text
            const helpText = document.createElement('small');
            helpText.className = 'form-text text-muted';
            helpText.innerHTML = 'System accounts must have a primary group';
            primaryGroupField.parentNode.appendChild(helpText);
        } else {
            // For other account types, primary group is optional
            primaryGroupField.removeAttribute('required');
            primaryGroupLabel.textContent = 'Primary Group';
            // Remove helper text if it exists
            const helpText = primaryGroupField.parentNode.querySelector('small');
            if (helpText) {
                helpText.remove();
            }
        }
        
        accountModal.show();
    }
    
    // Load groups for the primary group dropdown
    async function loadGroupsForDropdown() {
        try {
            // For system accounts, only system groups should be available as primary groups
            let groupsToLoad = [];
            if (sectionType === 'system') {
                // Only load system groups for system accounts
                groupsToLoad = await apiRequest(`/api/groups?type=system`);
            } else {
                // Load all groups of the current section type
                groupsToLoad = await apiRequest(`/api/groups?type=${sectionType}`);
            }
            
            const primaryGroupDropdown = document.getElementById('primaryGroupID');
            
            // Clear options except the first one (None)
            while (primaryGroupDropdown.options.length > 1) {
                primaryGroupDropdown.remove(1);
            }
            
            // Add groups to dropdown
            groupsToLoad.forEach(group => {
                const option = document.createElement('option');
                option.value = group.id;
                option.text = `${group.groupname} (GID: ${group.gid})`;
                primaryGroupDropdown.appendChild(option);
            });
            
            // For system accounts, select the first group by default if available
            if (sectionType === 'system' && groupsToLoad.length > 0) {
                primaryGroupDropdown.selectedIndex = 1; // Select first group
            }
        } catch (error) {
            console.error('Error loading groups for dropdown:', error);
            showError('Failed to load groups: ' + error.message);
        }
    }
    
    // Open modal for new group
    async function openNewGroupModal() {
        if (!groupModal) return;
        
        document.getElementById('groupModalLabel').textContent = 'New Group';
        document.getElementById('groupForm').reset();
        document.getElementById('groupId').value = '';
        
        // Clear any previous validation indicators
        const gidField = document.getElementById('gid');
        if (gidField) {
            gidField.classList.remove('is-valid', 'is-invalid');
            gidField.style.backgroundColor = '';
            const feedback = document.getElementById('gidFeedback');
            if (feedback) feedback.textContent = '';
            
            // Set min/max based on section type
            switch (sectionType) {
                case 'people':
                    gidField.setAttribute('min', '1000');
                    gidField.setAttribute('max', '60000');
                    gidField.setAttribute('placeholder', '1000-60000 (recommended)');
                    break;
                case 'system':
                    gidField.setAttribute('min', '9000');
                    gidField.setAttribute('max', '9100');
                    gidField.setAttribute('placeholder', '9000-9100 (recommended)');
                    break;
                case 'service':
                    gidField.setAttribute('min', '60001');
                    gidField.setAttribute('max', '65535');
                    gidField.setAttribute('placeholder', '60001-65535 (recommended)');
                    break;
                case 'database':
                    gidField.setAttribute('min', '70000');
                    gidField.setAttribute('max', '79999');
                    gidField.setAttribute('placeholder', '70000-79999 (recommended)');
                    break;
            }
        }
        
        // Fetch the next available GID and set it as default value
        try {
            const response = await authFetch(`/api/groups/next-gid?type=${sectionType}`);
            const data = await response.json();
            if (response.ok && data && data.gid) {
                gidField.value = data.gid;
                gidField.setAttribute('placeholder', `Suggestion: ${data.gid}`);
            }
        } catch (error) {
            console.warn('Error fetching next available GID:', error);
        }
        
        groupModal.show();
    }
    
    // Edit account
    async function editAccount(id) {
        if (!accountModal) return;
        
        try {
            // Load groups for the primary group dropdown first
            await loadGroupsForDropdown();
            
            const account = await apiRequest(`/api/accounts/${id}`);
            
            document.getElementById('accountModalLabel').textContent = 'Edit Account';
            document.getElementById('accountId').value = account.id;
            document.getElementById('uid').value = account.uid;
            document.getElementById('username').value = account.username;
            
            // Clear validation indicators
            const uidField = document.getElementById('uid');
            if (uidField) {
                uidField.classList.remove('is-valid', 'is-invalid');
                const feedback = document.getElementById('uidFeedback');
                if (feedback) feedback.textContent = '';
            }
            
            // Set primary group if it exists
            if (account.primary_group_id) {
                document.getElementById('primaryGroupID').value = account.primary_group_id;
            } else {
                document.getElementById('primaryGroupID').value = '';
            }
            
            // Verify account type matches section type
            const currentSectionType = document.getElementById('accountType').value;
            if (account.type.toLowerCase() !== currentSectionType.toLowerCase()) {
                console.warn(`Account type (${account.type}) doesn't match current section (${currentSectionType})`);
                showError(`Warning: This account is of type ${account.type}, but you're editing it from the ${currentSectionType} section`);
            }
            
            accountModal.show();
        } catch (error) {
            console.error('Error fetching account details:', error);
            showError('Failed to fetch account details: ' + error.message);
        }
    }
    
    // Edit group
    async function editGroup(id) {
        if (!groupModal) return;
        
        try {
            const group = await apiRequest(`/api/groups/${id}`);
            document.getElementById('groupModalLabel').textContent = 'Edit Group';
            document.getElementById('groupId').value = group.id;
            document.getElementById('gid').value = group.gid;
            document.getElementById('groupname').value = group.groupname;
            
            // Clear validation indicators
            const gidField = document.getElementById('gid');
            if (gidField) {
                gidField.classList.remove('is-valid', 'is-invalid');
                const feedback = document.getElementById('gidFeedback');
                if (feedback) feedback.textContent = '';
            }
            
            // Ensure description is never null or undefined
            document.getElementById('groupDescription').value = (group.description !== null && group.description !== undefined) ? group.description : '';
            
            // Set created_by field if it exists
            document.getElementById('createdBy').value = (group.created_by !== null && group.created_by !== undefined) ? group.created_by : '';
            
            // Verify group type matches section type
            if (group.type.toLowerCase() !== sectionType.toLowerCase()) {
                console.warn(`Group type (${group.type}) doesn't match current section (${sectionType})`);
                showError(`Warning: This group is of type ${group.type}, but you're editing it from the ${sectionType} section`);
            }
            
            groupModal.show();
        } catch (error) {
            console.error('Error fetching group details:', error);
            showError('Failed to fetch group details: ' + error.message);
        }
    }
    
    // Save account
    async function saveAccount() {
        try {
            // Validate UID first (now async)
            const isValidUID = await validateUID();
            if (!isValidUID) {
                return;
            }
            
            const accountId = document.getElementById('accountId').value;
            const uid = parseInt(document.getElementById('uid').value);
            const username = document.getElementById('username').value;
            const type = document.getElementById('accountType').value;
            const primaryGroupID = document.getElementById('primaryGroupID').value;
            
            // Validate username
            if (!username || username.trim() === '') {
                showError('Username is required');
                return;
            }
            
            // Validate primary group for system accounts
            if (type === 'system' && !primaryGroupID) {
                showError('System accounts must have a primary group');
                return;
            }
            
            const accountData = {
                uid,
                username,
                type,
                primary_group_id: primaryGroupID ? parseInt(primaryGroupID) : null
            };
            
            let result;
            if (accountId) {
                // Update existing account
                result = await apiRequest(`/api/accounts/${accountId}`, 'PUT', accountData);
                showSuccess('Account updated successfully');
            } else {
                // Create new account
                result = await apiRequest('/api/accounts', 'POST', accountData);
                showSuccess('Account created successfully');
            }
            
            if (accountModal) {
                accountModal.hide();
            }
            
            await loadAccounts();
        } catch (error) {
            console.error('Error saving account:', error);
            showError('Failed to save account: ' + (error.message || 'Unknown error'));
        }
    }
    
    // Save group
    async function saveGroup() {
        try {
            // Validate GID first (now async)
            const isValidGID = await validateGID();
            if (!isValidGID) {
                return;
            }
            
            const groupId = document.getElementById('groupId').value;
            const gid = parseInt(document.getElementById('gid').value);
            const groupname = document.getElementById('groupname').value;
            const description = document.getElementById('groupDescription').value;
            const createdBy = document.getElementById('createdBy').value;
            const type = document.getElementById('accountType').value;
            
            // Validate groupname
            if (!groupname || groupname.trim() === '') {
                showError('Groupname is required');
                return;
            }
            
            // Ensure description and created_by are never undefined or null
            const groupData = {
                gid,
                groupname,
                description: description || '',
                type,
                created_by: createdBy || ''
            };
            
            let result;
            if (groupId) {
                // Update existing group
                result = await apiRequest(`/api/groups/${groupId}`, 'PUT', groupData);
                showSuccess('Group updated successfully');
            } else {
                // Create new group
                result = await apiRequest('/api/groups', 'POST', groupData);
                showSuccess('Group created successfully');
            }
            
            if (groupModal) {
                groupModal.hide();
            }
            
            await loadGroups();
        } catch (error) {
            console.error('Error saving group:', error);
            showError('Failed to save group: ' + (error.message || 'Unknown error'));
        }
    }
    
    // Delete account
    async function deleteAccount(id) {
        if (!confirm('Are you sure you want to delete this account?')) {
            return;
        }
        
        try {
            await apiRequest(`/api/accounts/${id}`, 'DELETE');
            showSuccess('Account deleted successfully');
            await loadAccounts();
        } catch (error) {
            console.error('Error deleting account:', error);
            showError('Failed to delete account: ' + error.message);
        }
    }
    
    // Delete group
    async function deleteGroup(id) {
        if (!confirm('Are you sure you want to delete this group?')) {
            return;
        }
        
        try {
            await apiRequest(`/api/groups/${id}`, 'DELETE');
            showSuccess('Group deleted successfully');
            await loadGroups();
        } catch (error) {
            console.error('Error deleting group:', error);
            showError('Failed to delete group: ' + error.message);
        }
    }
    
    // Manage account groups (more advanced than just viewing)
    async function manageAccountGroups(id) {
        try {
            const [account, accountGroups, allGroups] = await Promise.all([
                apiRequest(`/api/accounts/${id}`),
                apiRequest(`/api/accounts/${id}/groups`),
                apiRequest(`/api/groups?type=${sectionType}`)
            ]);
            
            // Create a modal for managing groups if it doesn't exist
            let modalElement = document.getElementById('membershipModal');
            if (!modalElement) {
                modalElement = document.createElement('div');
                modalElement.id = 'membershipModal';
                modalElement.className = 'modal fade';
                modalElement.setAttribute('tabindex', '-1');
                modalElement.innerHTML = `
                    <div class="modal-dialog modal-lg">
                        <div class="modal-content">
                            <div class="modal-header">
                                <h5 class="modal-title">Manage Group Memberships</h5>
                                <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                            </div>
                            <div class="modal-body">
                                <div class="row">
                                    <div class="col-md-6">
                                        <h6>Account</h6>
                                        <p id="membershipAccountInfo"></p>
                                    </div>
                                    <div class="col-md-6">
                                        <h6>Available Groups</h6>
                                        <select id="availableGroups" class="form-select" size="8"></select>
                                    </div>
                                </div>
                                <div class="row mt-3">
                                    <div class="col-12">
                                        <h6>Current Group Memberships</h6>
                                        <ul id="currentGroups" class="list-group"></ul>
                                    </div>
                                </div>
                            </div>
                            <div class="modal-footer">
                                <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
                                <button type="button" class="btn btn-primary" id="addToGroupBtn">Add to Selected Group</button>
                            </div>
                        </div>
                    </div>
                `;
                document.body.appendChild(modalElement);
                membershipModal = new bootstrap.Modal(modalElement);
            }
            
            // Get elements
            const accountInfoElement = document.getElementById('membershipAccountInfo');
            const availableGroupsElement = document.getElementById('availableGroups');
            const currentGroupsElement = document.getElementById('currentGroups');
            const addToGroupBtn = document.getElementById('addToGroupBtn');
            
            // Set account info
            accountInfoElement.innerHTML = `
                <strong>${account.username}</strong> (UID: ${account.uid})<br>
                Type: ${account.type}<br>
                ID: ${account.id}
            `;
            
            // Clear available groups dropdown
            availableGroupsElement.innerHTML = '';
            
            // Populate available groups (filtering out groups the account is already in)
            const accountGroupIds = accountGroups.map(g => g.id);
            const availableGroupsFiltered = allGroups.filter(g => !accountGroupIds.includes(g.id));
            
            if (availableGroupsFiltered.length === 0) {
                const option = document.createElement('option');
                option.disabled = true;
                option.textContent = 'No more available groups';
                availableGroupsElement.appendChild(option);
            } else {
                availableGroupsFiltered.forEach(group => {
                    const option = document.createElement('option');
                    option.value = group.id;
                    option.textContent = `${group.groupname} (GID: ${group.gid})`;
                    availableGroupsElement.appendChild(option);
                });
            }
            
            // Populate current groups
            currentGroupsElement.innerHTML = '';
            if (accountGroups.length === 0) {
                const li = document.createElement('li');
                li.className = 'list-group-item text-center';
                li.textContent = 'Not a member of any groups';
                currentGroupsElement.appendChild(li);
            } else {
                accountGroups.forEach(group => {
                    const li = document.createElement('li');
                    li.className = 'list-group-item d-flex justify-content-between align-items-center';
                    li.innerHTML = `
                        ${group.groupname} <span class="badge bg-primary">GID: ${group.gid}</span>
                        <button type="button" class="btn btn-sm btn-danger remove-from-group" data-group-id="${group.id}">Remove</button>
                    `;
                    currentGroupsElement.appendChild(li);
                });
                
                // Add event listeners to remove buttons
                document.querySelectorAll('.remove-from-group').forEach(btn => {
                    btn.addEventListener('click', async (e) => {
                        const groupId = e.target.dataset.groupId;
                        try {
                            await apiRequest('/api/memberships', 'DELETE', {
                                account_id: account.id,
                                group_id: parseInt(groupId)
                            });
                            // Refresh the modal content
                            manageAccountGroups(account.id);
                        } catch (error) {
                            console.error('Error removing from group:', error);
                            showError('Failed to remove from group: ' + error.message);
                        }
                    });
                });
            }
            
            // Handle add to group button
            addToGroupBtn.onclick = async () => {
                const selectedGroup = availableGroupsElement.value;
                if (!selectedGroup) {
                    showError('Please select a group to add');
                    return;
                }
                
                try {
                    await apiRequest('/api/memberships', 'POST', {
                        account_id: account.id,
                        group_id: parseInt(selectedGroup)
                    });
                    // Refresh the modal content
                    manageAccountGroups(account.id);
                } catch (error) {
                    console.error('Error adding to group:', error);
                    showError('Failed to add to group: ' + error.message);
                }
            };
            
            // Show the modal
            membershipModal.show();
        } catch (error) {
            console.error('Error fetching membership data:', error);
            showError('Failed to fetch membership data: ' + error.message);
        }
    }
    
    // Manage group members (more advanced than just viewing)
    async function manageGroupMembers(id) {
        try {
            const [group, groupMembers, allAccounts] = await Promise.all([
                apiRequest(`/api/groups/${id}`),
                apiRequest(`/api/groups/${id}/accounts`),
                apiRequest(`/api/accounts?type=${sectionType}`)
            ]);
            
            // Create a modal for managing members if it doesn't exist
            let modalElement = document.getElementById('membershipModal');
            if (!modalElement) {
                modalElement = document.createElement('div');
                modalElement.id = 'membershipModal';
                modalElement.className = 'modal fade';
                modalElement.setAttribute('tabindex', '-1');
                modalElement.innerHTML = `
                    <div class="modal-dialog modal-lg">
                        <div class="modal-content">
                            <div class="modal-header">
                                <h5 class="modal-title">Manage Group Members</h5>
                                <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                            </div>
                            <div class="modal-body">
                                <div class="row">
                                    <div class="col-md-6">
                                        <h6>Group</h6>
                                        <p id="membershipGroupInfo"></p>
                                    </div>
                                    <div class="col-md-6">
                                        <h6>Available Accounts</h6>
                                        <select id="availableAccounts" class="form-select" size="8"></select>
                                    </div>
                                </div>
                                <div class="row mt-3">
                                    <div class="col-12">
                                        <h6>Current Members</h6>
                                        <ul id="currentMembers" class="list-group"></ul>
                                    </div>
                                </div>
                            </div>
                            <div class="modal-footer">
                                <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
                                <button type="button" class="btn btn-primary" id="addMemberBtn">Add Selected Account</button>
                            </div>
                        </div>
                    </div>
                `;
                document.body.appendChild(modalElement);
                membershipModal = new bootstrap.Modal(modalElement);
            }
            
            // Get elements
            const groupInfoElement = document.getElementById('membershipGroupInfo');
            const availableAccountsElement = document.getElementById('availableAccounts');
            const currentMembersElement = document.getElementById('currentMembers');
            const addMemberBtn = document.getElementById('addMemberBtn');
            
            // Set group info
            groupInfoElement.innerHTML = `
                <strong>${group.groupname}</strong> (GID: ${group.gid})<br>
                Type: ${group.type}<br>
                ID: ${group.id}
            `;
            
            // Clear available accounts dropdown
            availableAccountsElement.innerHTML = '';
            
            // Populate available accounts (filtering out accounts already in the group)
            const memberIds = groupMembers.map(a => a.id);
            const availableAccountsFiltered = allAccounts.filter(a => !memberIds.includes(a.id));
            
            if (availableAccountsFiltered.length === 0) {
                const option = document.createElement('option');
                option.disabled = true;
                option.textContent = 'No more available accounts';
                availableAccountsElement.appendChild(option);
            } else {
                availableAccountsFiltered.forEach(account => {
                    const option = document.createElement('option');
                    option.value = account.id;
                    option.textContent = `${account.username} (UID: ${account.uid})`;
                    availableAccountsElement.appendChild(option);
                });
            }
            
            // Populate current members
            currentMembersElement.innerHTML = '';
            if (groupMembers.length === 0) {
                const li = document.createElement('li');
                li.className = 'list-group-item text-center';
                li.textContent = 'No members in this group';
                currentMembersElement.appendChild(li);
            } else {
                groupMembers.forEach(account => {
                    const li = document.createElement('li');
                    li.className = 'list-group-item d-flex justify-content-between align-items-center';
                    li.innerHTML = `
                        ${account.username} <span class="badge bg-primary">UID: ${account.uid}</span>
                        <button type="button" class="btn btn-sm btn-danger remove-member" data-account-id="${account.id}">Remove</button>
                    `;
                    currentMembersElement.appendChild(li);
                });
                
                // Add event listeners to remove buttons
                document.querySelectorAll('.remove-member').forEach(btn => {
                    btn.addEventListener('click', async (e) => {
                        const accountId = e.target.dataset.accountId;
                        try {
                            await apiRequest('/api/memberships', 'DELETE', {
                                account_id: parseInt(accountId),
                                group_id: group.id
                            });
                            // Refresh the modal content
                            manageGroupMembers(group.id);
                        } catch (error) {
                            console.error('Error removing member:', error);
                            showError('Failed to remove member: ' + error.message);
                        }
                    });
                });
            }
            
            // Handle add member button
            addMemberBtn.onclick = async () => {
                const selectedAccount = availableAccountsElement.value;
                if (!selectedAccount) {
                    showError('Please select an account to add');
                    return;
                }
                
                try {
                    await apiRequest('/api/memberships', 'POST', {
                        account_id: parseInt(selectedAccount),
                        group_id: group.id
                    });
                    // Refresh the modal content
                    manageGroupMembers(group.id);
                } catch (error) {
                    console.error('Error adding member:', error);
                    showError('Failed to add member: ' + error.message);
                }
            };
            
            // Show the modal
            membershipModal.show();
        } catch (error) {
            console.error('Error fetching membership data:', error);
            showError('Failed to fetch membership data: ' + error.message);
        }
    }
    
    // View account groups (simple view)
    async function viewAccountGroups(id) {
        try {
            const groups = await apiRequest(`/api/accounts/${id}/groups`);
            if (groups.length === 0) {
                showInfo('This account is not a member of any groups');
                return;
            }
            
            let message = '<h5>Group Memberships</h5><ul>';
            groups.forEach(group => {
                message += `<li>${group.groupname} (GID: ${group.gid})</li>`;
            });
            message += '</ul>';
            
            showModal('Account Group Memberships', message);
        } catch (error) {
            console.error('Error fetching account groups:', error);
            showError('Failed to fetch account groups: ' + error.message);
        }
    }
    
    // View group members (simple view)
    async function viewGroupMembers(id) {
        try {
            const accounts = await apiRequest(`/api/groups/${id}/accounts`);
            if (accounts.length === 0) {
                showInfo('This group has no members');
                return;
            }
            
            let message = '<h5>Group Members</h5><ul>';
            accounts.forEach(account => {
                message += `<li>${account.username} (UID: ${account.uid})</li>`;
            });
            message += '</ul>';
            
            showModal('Group Members', message);
        } catch (error) {
            console.error('Error fetching group members:', error);
            showError('Failed to fetch group members: ' + error.message);
        }
    }
    
    // Helper function to show a modal with content
    function showModal(title, content) {
        // Create a modal if it doesn't exist
        let infoModal = document.getElementById('infoModal');
        if (!infoModal) {
            infoModal = document.createElement('div');
            infoModal.id = 'infoModal';
            infoModal.className = 'modal fade';
            infoModal.setAttribute('tabindex', '-1');
            infoModal.innerHTML = `
                <div class="modal-dialog">
                    <div class="modal-content">
                        <div class="modal-header">
                            <h5 class="modal-title" id="infoModalTitle"></h5>
                            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                        </div>
                        <div class="modal-body" id="infoModalContent"></div>
                        <div class="modal-footer">
                            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
                        </div>
                    </div>
                </div>
            `;
            document.body.appendChild(infoModal);
        }
        
        // Set content
        document.getElementById('infoModalTitle').textContent = title;
        document.getElementById('infoModalContent').innerHTML = content;
        
        // Show the modal
        const modal = new bootstrap.Modal(infoModal);
        modal.show();
    }
    
    // Helper function to show info message
    function showInfo(message) {
        const infoDiv = document.createElement('div');
        infoDiv.className = 'alert alert-info alert-dismissible fade show';
        infoDiv.innerHTML = `
            ${message}
            <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
        `;
        
        // Insert at the top of the container
        const container = document.querySelector('.container');
        if (container) {
            container.insertBefore(infoDiv, container.firstChild);
            
            // Auto-dismiss after 5 seconds
            setTimeout(() => {
                infoDiv.remove();
            }, 5000);
        }
    }
    
    // Load audit logs with pagination and filtering
    async function loadAuditLogs(page = 1) {
        if (!auditTableBody) return;
        
        auditCurrentPage = page;
        auditTableBody.innerHTML = '<tr><td colspan="7" class="text-center">Loading audit logs...</td></tr>';
        
        try {
            // Get filter values
            const entityType = auditFilterEntity ? auditFilterEntity.value : '';
            const actionType = auditFilterAction ? auditFilterAction.value : '';
            
            // Build query parameters
            let queryParams = `page=${page}&limit=${auditPageSize}`;
            if (entityType) queryParams += `&entity_type=${entityType}`;
            if (actionType) queryParams += `&action=${actionType}`;
            if (sectionType) queryParams += `&section=${sectionType}`;
            
            // Add a timestamp parameter to avoid caching issues
            const timestamp = new Date().getTime();
            queryParams += `&_=${timestamp}`;
            
            // Get audit logs from API
            const result = await apiRequest(`/api/audit?${queryParams}`);
            
            // Update pagination information
            auditTotalPages = Math.ceil(result.total / auditPageSize);
            
            // Render audit logs
            renderAuditTable(result.entries);
            
            // Update pagination UI
            updateAuditPagination();
            
        } catch (error) {
            console.error('Error loading audit logs:', error);
            showError('Failed to load audit logs: ' + error.message);
            
            // Show empty state
            auditTableBody.innerHTML = '<tr><td colspan="7" class="text-center text-danger">Error loading audit logs</td></tr>';
        }
    }
    
    // Render audit logs table
    function renderAuditTable(auditLogs) {
        if (!auditTableBody) return;
        
        auditTableBody.innerHTML = '';
        
        if (!auditLogs || auditLogs.length === 0) {
            const row = document.createElement('tr');
            row.innerHTML = '<td colspan="7" class="text-center">No audit logs found</td>';
            auditTableBody.appendChild(row);
            return;
        }
        
        auditLogs.forEach(entry => {
            // Get badge color based on action
            let badgeClass = 'bg-secondary';
            switch (entry.action) {
                case 'create':
                    badgeClass = 'bg-success';
                    break;
                case 'update':
                    badgeClass = 'bg-primary';
                    break;
                case 'delete':
                    badgeClass = 'bg-danger';
                    break;
                case 'assign':
                    badgeClass = 'bg-info';
                    break;
                case 'remove':
                    badgeClass = 'bg-warning';
                    break;
            }
            
            // Format details for display (truncate if too long)
            let details = entry.details;
            if (details && details.length > 50) {
                details = details.substring(0, 50) + '...';
            }
            
            // Create table row
            const row = document.createElement('tr');
            row.className = 'audit-row';
            
            // Format entity type with icon
            let entityIcon = 'bi-question-circle';
            switch(entry.entity_type) {
                case 'account':
                    entityIcon = 'bi-person';
                    break;
                case 'group':
                    entityIcon = 'bi-people';
                    break;
                case 'account_group':
                    entityIcon = 'bi-person-plus';
                    break;
            }
            
            // Get action icon
            let actionIcon = 'bi-activity';
            switch(entry.action) {
                case 'create':
                    actionIcon = 'bi-plus-circle';
                    break;
                case 'update':
                    actionIcon = 'bi-pencil';
                    break;
                case 'delete':
                    actionIcon = 'bi-trash';
                    break;
                case 'assign':
                    actionIcon = 'bi-link';
                    break;
                case 'remove':
                    actionIcon = 'bi-unlink';
                    break;
            }
            
            row.innerHTML = `
                <td>${entry.id}</td>
                <td><span title="${new Date(entry.created_at).toLocaleString()}">${formatDate(entry.created_at)}</span></td>
                <td><span class="badge ${badgeClass}"><i class="bi ${actionIcon} me-1"></i>${entry.action}</span></td>
                <td><i class="bi ${entityIcon} me-1"></i>${entry.entity_type}</td>
                <td>${entry.entity_id}</td>
                <td><i class="bi bi-person-circle me-1"></i>${entry.username || 'System'}</td>
                <td class="text-end">
                    <button class="btn btn-sm btn-outline-primary view-audit rounded-pill" data-id="${entry.id}">
                        <i class="bi bi-eye me-1"></i> View Details
                    </button>
                </td>
            `;
            auditTableBody.appendChild(row);
        });
        
        // Add event listeners to view buttons
        document.querySelectorAll('.view-audit').forEach(btn => {
            btn.addEventListener('click', () => viewAuditDetail(btn.dataset.id));
        });
    }
    
    // View audit detail
    async function viewAuditDetail(id) {
        if (!auditDetailModal) return;
        
        try {
            // Get audit entry details
            const entry = await apiRequest(`/api/audit/${id}`);
            
            // Set modal content
            document.getElementById('auditDetailId').textContent = entry.id;
            document.getElementById('auditDetailTimestamp').textContent = formatDate(entry.created_at);
            
            // Set badge for action
            const actionBadge = document.getElementById('auditDetailAction');
            actionBadge.textContent = entry.action;
            
            // Set badge color based on action
            let badgeClass = 'bg-secondary';
            switch (entry.action) {
                case 'create':
                    badgeClass = 'bg-success';
                    break;
                case 'update':
                    badgeClass = 'bg-primary';
                    break;
                case 'delete':
                    badgeClass = 'bg-danger';
                    break;
                case 'assign':
                    badgeClass = 'bg-info';
                    break;
                case 'remove':
                    badgeClass = 'bg-warning';
                    break;
            }
            actionBadge.classList.remove('bg-success', 'bg-primary', 'bg-danger', 'bg-info', 'bg-warning', 'bg-secondary');
            actionBadge.classList.add(badgeClass);
            
            document.getElementById('auditDetailEntityType').textContent = entry.entity_type;
            document.getElementById('auditDetailEntityId').textContent = entry.entity_id;
            document.getElementById('auditDetailUser').textContent = entry.username || 'System';
            document.getElementById('auditDetailIP').textContent = entry.ip_address || 'N/A';
            
            // Format details as JSON if it's a JSON string
            let formattedDetails = entry.details || '';
            try {
                if (formattedDetails && typeof formattedDetails === 'string' && (formattedDetails.startsWith('{') || formattedDetails.startsWith('['))) {
                    formattedDetails = JSON.stringify(JSON.parse(formattedDetails), null, 2);
                }
            } catch (e) {
                console.warn('Could not parse details as JSON:', e);
            }
            
            document.getElementById('auditDetailContent').textContent = formattedDetails;
            
            // Show the modal
            auditDetailModal.show();
        } catch (error) {
            console.error('Error fetching audit details:', error);
            showError('Failed to fetch audit details: ' + error.message);
        }
    }
    
    // Update audit pagination UI
    function updateAuditPagination() {
        if (!auditPagination) return;
        
        auditPagination.innerHTML = '';
        
        // If no pages, don't show pagination
        if (auditTotalPages <= 1) {
            return;
        }
        
        // Information text about pages
        const infoLi = document.createElement('li');
        infoLi.className = 'page-item disabled d-none d-md-block';
        infoLi.innerHTML = `
            <span class="page-link border-0 bg-transparent">
                Page ${auditCurrentPage} of ${auditTotalPages}
            </span>
        `;
        auditPagination.appendChild(infoLi);
        
        // First page button (only show if not near first page)
        if (auditCurrentPage > 3) {
            const firstLi = document.createElement('li');
            firstLi.className = 'page-item d-none d-md-block';
            firstLi.innerHTML = `
                <button class="page-link" aria-label="First Page">
                    <i class="bi bi-chevron-double-left"></i>
                </button>
            `;
            firstLi.querySelector('button').addEventListener('click', () => loadAuditLogs(1));
            auditPagination.appendChild(firstLi);
        }
        
        // Previous button
        const prevLi = document.createElement('li');
        prevLi.className = `page-item ${auditCurrentPage === 1 ? 'disabled' : ''}`;
        prevLi.innerHTML = `
            <button class="page-link" ${auditCurrentPage === 1 ? 'disabled' : ''} aria-label="Previous">
                <i class="bi bi-chevron-left"></i>
            </button>
        `;
        auditPagination.appendChild(prevLi);
        
        // Only load button if not disabled
        if (auditCurrentPage !== 1) {
            prevLi.querySelector('button').addEventListener('click', () => loadAuditLogs(auditCurrentPage - 1));
        }
        
        // Page buttons - show fewer on mobile
        const isMobile = window.innerWidth < 768;
        const pagesDisplayed = isMobile ? 3 : 5;
        const offset = Math.floor(pagesDisplayed / 2);
        
        let startPage = Math.max(1, auditCurrentPage - offset);
        let endPage = Math.min(auditTotalPages, startPage + pagesDisplayed - 1);
        
        // Adjust if we're near the end
        if (endPage - startPage + 1 < pagesDisplayed) {
            startPage = Math.max(1, endPage - pagesDisplayed + 1);
        }
        
        // Add page buttons
        for (let i = startPage; i <= endPage; i++) {
            const pageLi = document.createElement('li');
            pageLi.className = `page-item ${i === auditCurrentPage ? 'active' : ''}`;
            pageLi.innerHTML = `
                <button class="page-link">${i}</button>
            `;
            auditPagination.appendChild(pageLi);
            
            // Add event listener to page button
            if (i !== auditCurrentPage) {
                pageLi.querySelector('button').addEventListener('click', () => loadAuditLogs(i));
            }
        }
        
        // Next button
        const nextLi = document.createElement('li');
        nextLi.className = `page-item ${auditCurrentPage === auditTotalPages ? 'disabled' : ''}`;
        nextLi.innerHTML = `
            <button class="page-link" ${auditCurrentPage === auditTotalPages ? 'disabled' : ''} aria-label="Next">
                <i class="bi bi-chevron-right"></i>
            </button>
        `;
        auditPagination.appendChild(nextLi);
        
        // Only load button if not disabled
        if (auditCurrentPage !== auditTotalPages) {
            nextLi.querySelector('button').addEventListener('click', () => loadAuditLogs(auditCurrentPage + 1));
        }
        
        // Last page button (only show if not near last page)
        if (auditCurrentPage < auditTotalPages - 2) {
            const lastLi = document.createElement('li');
            lastLi.className = 'page-item d-none d-md-block';
            lastLi.innerHTML = `
                <button class="page-link" aria-label="Last Page">
                    <i class="bi bi-chevron-double-right"></i>
                </button>
            `;
            lastLi.querySelector('button').addEventListener('click', () => loadAuditLogs(auditTotalPages));
            auditPagination.appendChild(lastLi);
        }
    }
});
