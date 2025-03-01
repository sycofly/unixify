// JavaScript for the section pages (people, system, database, service)
document.addEventListener('DOMContentLoaded', function() {
    // Get section type from the hidden input
    const sectionType = document.getElementById('accountType').value;
    
    // Elements
    const accountsTableBody = document.getElementById('accountsTableBody');
    const groupsTableBody = document.getElementById('groupsTableBody');
    const newAccountBtn = document.getElementById('newAccountBtn');
    const newGroupBtn = document.getElementById('newGroupBtn');
    const saveAccountBtn = document.getElementById('saveAccountBtn');
    const saveGroupBtn = document.getElementById('saveGroupBtn');
    
    // Bootstrap modals
    let accountModal, groupModal;
    if (document.getElementById('accountModal')) {
        accountModal = new bootstrap.Modal(document.getElementById('accountModal'));
    }
    if (document.getElementById('groupModal')) {
        groupModal = new bootstrap.Modal(document.getElementById('groupModal'));
    }
    
    // Load accounts and groups on page load
    loadAccounts();
    loadGroups();
    
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
    
    // Functions
    
    // Load accounts for the current section
    async function loadAccounts() {
        try {
            console.log(`Loading accounts for section type: ${sectionType}`);
            // Add a timestamp parameter to avoid caching issues
            const timestamp = new Date().getTime();
            const accounts = await apiRequest(`/api/accounts?type=${sectionType}&_=${timestamp}`);
            console.log('Accounts loaded:', accounts);
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
            console.log(`Loading groups for section type: ${sectionType}`);
            // Add a timestamp parameter to avoid caching issues
            const timestamp = new Date().getTime();
            const groups = await apiRequest(`/api/groups?type=${sectionType}&_=${timestamp}`);
            console.log('Groups loaded:', groups);
            renderGroupsTable(groups);
            return groups;
        } catch (error) {
            console.error('Error loading groups:', error);
            showError('Failed to load groups: ' + error.message);
            return [];
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
            console.log('Rendering account:', account);
            
            // Get primary group info - with improved handling
            let primaryGroupText = 'None';
            if (account.primary_group && account.primary_group.groupname) {
                primaryGroupText = `${account.primary_group.groupname} (${account.primary_group.gid})`;
            } else if (account.primary_group_id) {
                // If we only have the ID but not the full group data, show the ID
                // This should not happen if the API correctly preloads the primary group
                primaryGroupText = `ID: ${account.primary_group_id}`;
                console.warn(`Account ${account.id} has primary_group_id ${account.primary_group_id} but no primary_group data`);
            }
            
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${account.id}</td>
                <td>${account.uid}</td>
                <td>${account.username}</td>
                <td>${primaryGroupText}</td>
                <td>${formatDate(account.created_at)}</td>
                <td class="action-buttons">
                    <button type="button" class="btn btn-sm btn-primary edit-account" data-id="${account.id}">Edit</button>
                    <button type="button" class="btn btn-sm btn-danger delete-account" data-id="${account.id}">Delete</button>
                    <button type="button" class="btn btn-sm btn-secondary groups-account" data-id="${account.id}">Groups</button>
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
            btn.addEventListener('click', () => viewAccountGroups(btn.dataset.id));
        });
    }
    
    // Render groups table
    function renderGroupsTable(groups) {
        if (!groupsTableBody) return;
        
        groupsTableBody.innerHTML = '';
        
        if (!groups || groups.length === 0) {
            const row = document.createElement('tr');
            row.innerHTML = '<td colspan="5" class="text-center">No groups found</td>';
            groupsTableBody.appendChild(row);
            return;
        }
        
        groups.forEach(group => {
            const row = document.createElement('tr');
            // Ensure description is displayed properly
            const displayDescription = group.description !== null && group.description !== undefined ? group.description : '';
            
            row.innerHTML = `
                <td>${group.id}</td>
                <td>${group.gid}</td>
                <td>${group.groupname}</td>
                <td>${displayDescription}</td>
                <td>${formatDate(group.created_at)}</td>
                <td class="action-buttons">
                    <button type="button" class="btn btn-sm btn-primary edit-group" data-id="${group.id}">Edit</button>
                    <button type="button" class="btn btn-sm btn-danger delete-group" data-id="${group.id}">Delete</button>
                    <button type="button" class="btn btn-sm btn-secondary members-group" data-id="${group.id}">Members</button>
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
            btn.addEventListener('click', () => viewGroupMembers(btn.dataset.id));
        });
    }
    
    // Open modal for new account
    function openNewAccountModal() {
        if (!accountModal) return;
        
        document.getElementById('accountModalLabel').textContent = 'New Account';
        document.getElementById('accountForm').reset();
        document.getElementById('accountId').value = '';
        
        // Load groups for the primary group dropdown
        loadGroupsForDropdown();
        
        // Set up form requirements based on section type
        const currentSectionType = document.getElementById('accountType').value;
        const primaryGroupField = document.getElementById('primaryGroupID');
        const primaryGroupLabel = primaryGroupField.closest('.mb-3').querySelector('.form-label');
        
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
            console.log(`Loading groups for dropdown with section type: ${sectionType}`);
            
            // For system accounts, only system groups should be available as primary groups
            let groupsToLoad = [];
            if (sectionType === 'system') {
                // Only load system groups for system accounts
                groupsToLoad = await apiRequest(`/api/groups?type=system`);
                console.log('System groups loaded for system account:', groupsToLoad);
            } else {
                // Load all groups of the current section type
                groupsToLoad = await apiRequest(`/api/groups?type=${sectionType}`);
                console.log(`${sectionType} groups loaded:`, groupsToLoad);
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
    function openNewGroupModal() {
        if (!groupModal) return;
        
        document.getElementById('groupModalLabel').textContent = 'New Group';
        document.getElementById('groupForm').reset();
        document.getElementById('groupId').value = '';
        groupModal.show();
    }
    
    // Edit account
    async function editAccount(id) {
        if (!accountModal) return;
        
        try {
            console.log(`Editing account with ID: ${id}`);
            
            // Load groups for the primary group dropdown first
            await loadGroupsForDropdown();
            
            const account = await apiRequest(`/api/accounts/${id}`);
            console.log('Account details fetched:', account);
            
            document.getElementById('accountModalLabel').textContent = 'Edit Account';
            document.getElementById('accountId').value = account.id;
            document.getElementById('uid').value = account.uid;
            document.getElementById('username').value = account.username;
            
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
            }
            
            console.log(`Account type: ${account.type}, Section type: ${currentSectionType}`);
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
            // Ensure description is never null or undefined
            document.getElementById('groupDescription').value = (group.description !== null && group.description !== undefined) ? group.description : '';
            groupModal.show();
        } catch (error) {
            console.error('Error fetching group details:', error);
            showError('Failed to fetch group details: ' + error.message);
        }
    }
    
    // Save account
    async function saveAccount() {
        console.log('Save account button clicked');
        
        try {
            const accountId = document.getElementById('accountId').value;
            const uid = parseInt(document.getElementById('uid').value);
            const username = document.getElementById('username').value;
            const type = document.getElementById('accountType').value;
            const primaryGroupID = document.getElementById('primaryGroupID').value;
            
            console.log('Account data collected:', { 
                accountId, 
                uid, 
                username, 
                type, 
                primaryGroupID 
            });
            
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
            
            console.log('Making API request with data:', accountData);
            
            let result;
            if (accountId) {
                // Update existing account
                console.log('Updating existing account with ID:', accountId);
                result = await apiRequest(`/api/accounts/${accountId}`, 'PUT', accountData);
                console.log('Account update response:', result);
                showSuccess('Account updated successfully');
            } else {
                // Create new account
                console.log('Creating new account');
                result = await apiRequest('/api/accounts', 'POST', accountData);
                console.log('API response:', result);
                showSuccess('Account created successfully');
            }
            
            if (accountModal) {
                console.log('Hiding modal');
                accountModal.hide();
            } else {
                console.log('Account modal not found');
            }
            
            console.log('Reloading accounts');
            await loadAccounts();
            console.log('Accounts reloaded');
        } catch (error) {
            console.error('Error saving account:', error);
            showError('Failed to save account: ' + (error.message || 'Unknown error'));
        }
    }
    
    // Save group
    async function saveGroup() {
        console.log('Save group button clicked');
        
        try {
            const groupId = document.getElementById('groupId').value;
            const gid = parseInt(document.getElementById('gid').value);
            const groupname = document.getElementById('groupname').value;
            const description = document.getElementById('groupDescription').value;
            const type = document.getElementById('accountType').value;
            
            console.log('Group data collected:', { 
                groupId, 
                gid, 
                groupname, 
                description, 
                type 
            });
            
            // Ensure description is never undefined or null
            const groupData = {
                gid,
                groupname,
                description: description || '',
                type
            };
            
            console.log('Making API request with data:', groupData);
            
            let result;
            if (groupId) {
                // Update existing group
                console.log('Updating existing group with ID:', groupId);
                result = await apiRequest(`/api/groups/${groupId}`, 'PUT', groupData);
                console.log('Group update response:', result);
                showSuccess('Group updated successfully');
            } else {
                // Create new group
                console.log('Creating new group');
                result = await apiRequest('/api/groups', 'POST', groupData);
                console.log('API response:', result);
                showSuccess('Group created successfully');
            }
            
            if (groupModal) {
                console.log('Hiding modal');
                groupModal.hide();
            } else {
                console.log('Group modal not found');
            }
            
            console.log('Reloading groups');
            loadGroups();
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
            loadAccounts();
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
            loadGroups();
        } catch (error) {
            console.error('Error deleting group:', error);
            showError('Failed to delete group: ' + error.message);
        }
    }
    
    // View account groups
    async function viewAccountGroups(id) {
        try {
            const groups = await apiRequest(`/api/accounts/${id}/groups`);
            const groupNames = groups.map(g => g.groupname).join(', ');
            alert(`This account is a member of ${groups.length} groups: ${groupNames || 'None'}`);
        } catch (error) {
            console.error('Error fetching account groups:', error);
            showError('Failed to fetch account groups: ' + error.message);
        }
    }
    
    // View group members
    async function viewGroupMembers(id) {
        try {
            const accounts = await apiRequest(`/api/groups/${id}/accounts`);
            const accountNames = accounts.map(a => a.username).join(', ');
            alert(`This group has ${accounts.length} members: ${accountNames || 'None'}`);
        } catch (error) {
            console.error('Error fetching group members:', error);
            showError('Failed to fetch group members: ' + error.message);
        }
    }
});
