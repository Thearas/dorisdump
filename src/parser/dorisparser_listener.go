// Code generated from DorisParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // DorisParser
import "github.com/antlr4-go/antlr/v4"

// DorisParserListener is a complete listener for a parse tree produced by DorisParser.
type DorisParserListener interface {
	antlr.ParseTreeListener

	// EnterMultiStatements is called when entering the multiStatements production.
	EnterMultiStatements(c *MultiStatementsContext)

	// EnterSingleStatement is called when entering the singleStatement production.
	EnterSingleStatement(c *SingleStatementContext)

	// EnterStatementBaseAlias is called when entering the statementBaseAlias production.
	EnterStatementBaseAlias(c *StatementBaseAliasContext)

	// EnterCallProcedure is called when entering the callProcedure production.
	EnterCallProcedure(c *CallProcedureContext)

	// EnterCreateProcedure is called when entering the createProcedure production.
	EnterCreateProcedure(c *CreateProcedureContext)

	// EnterDropProcedure is called when entering the dropProcedure production.
	EnterDropProcedure(c *DropProcedureContext)

	// EnterShowProcedureStatus is called when entering the showProcedureStatus production.
	EnterShowProcedureStatus(c *ShowProcedureStatusContext)

	// EnterShowCreateProcedure is called when entering the showCreateProcedure production.
	EnterShowCreateProcedure(c *ShowCreateProcedureContext)

	// EnterShowConfig is called when entering the showConfig production.
	EnterShowConfig(c *ShowConfigContext)

	// EnterStatementDefault is called when entering the statementDefault production.
	EnterStatementDefault(c *StatementDefaultContext)

	// EnterSupportedDmlStatementAlias is called when entering the supportedDmlStatementAlias production.
	EnterSupportedDmlStatementAlias(c *SupportedDmlStatementAliasContext)

	// EnterSupportedCreateStatementAlias is called when entering the supportedCreateStatementAlias production.
	EnterSupportedCreateStatementAlias(c *SupportedCreateStatementAliasContext)

	// EnterSupportedAlterStatementAlias is called when entering the supportedAlterStatementAlias production.
	EnterSupportedAlterStatementAlias(c *SupportedAlterStatementAliasContext)

	// EnterMaterailizedViewStatementAlias is called when entering the materailizedViewStatementAlias production.
	EnterMaterailizedViewStatementAlias(c *MaterailizedViewStatementAliasContext)

	// EnterConstraintStatementAlias is called when entering the constraintStatementAlias production.
	EnterConstraintStatementAlias(c *ConstraintStatementAliasContext)

	// EnterSupportedDropStatementAlias is called when entering the supportedDropStatementAlias production.
	EnterSupportedDropStatementAlias(c *SupportedDropStatementAliasContext)

	// EnterUnsupported is called when entering the unsupported production.
	EnterUnsupported(c *UnsupportedContext)

	// EnterUnsupportedStatement is called when entering the unsupportedStatement production.
	EnterUnsupportedStatement(c *UnsupportedStatementContext)

	// EnterCreateMTMV is called when entering the createMTMV production.
	EnterCreateMTMV(c *CreateMTMVContext)

	// EnterRefreshMTMV is called when entering the refreshMTMV production.
	EnterRefreshMTMV(c *RefreshMTMVContext)

	// EnterAlterMTMV is called when entering the alterMTMV production.
	EnterAlterMTMV(c *AlterMTMVContext)

	// EnterDropMTMV is called when entering the dropMTMV production.
	EnterDropMTMV(c *DropMTMVContext)

	// EnterPauseMTMV is called when entering the pauseMTMV production.
	EnterPauseMTMV(c *PauseMTMVContext)

	// EnterResumeMTMV is called when entering the resumeMTMV production.
	EnterResumeMTMV(c *ResumeMTMVContext)

	// EnterCancelMTMVTask is called when entering the cancelMTMVTask production.
	EnterCancelMTMVTask(c *CancelMTMVTaskContext)

	// EnterShowCreateMTMV is called when entering the showCreateMTMV production.
	EnterShowCreateMTMV(c *ShowCreateMTMVContext)

	// EnterAddConstraint is called when entering the addConstraint production.
	EnterAddConstraint(c *AddConstraintContext)

	// EnterDropConstraint is called when entering the dropConstraint production.
	EnterDropConstraint(c *DropConstraintContext)

	// EnterShowConstraint is called when entering the showConstraint production.
	EnterShowConstraint(c *ShowConstraintContext)

	// EnterInsertTable is called when entering the insertTable production.
	EnterInsertTable(c *InsertTableContext)

	// EnterUpdate is called when entering the update production.
	EnterUpdate(c *UpdateContext)

	// EnterDelete is called when entering the delete production.
	EnterDelete(c *DeleteContext)

	// EnterLoad is called when entering the load production.
	EnterLoad(c *LoadContext)

	// EnterMysqlLoad is called when entering the mysqlLoad production.
	EnterMysqlLoad(c *MysqlLoadContext)

	// EnterExport is called when entering the export production.
	EnterExport(c *ExportContext)

	// EnterCreateTable is called when entering the createTable production.
	EnterCreateTable(c *CreateTableContext)

	// EnterCreateView is called when entering the createView production.
	EnterCreateView(c *CreateViewContext)

	// EnterCreateTableLike is called when entering the createTableLike production.
	EnterCreateTableLike(c *CreateTableLikeContext)

	// EnterCreateRowPolicy is called when entering the createRowPolicy production.
	EnterCreateRowPolicy(c *CreateRowPolicyContext)

	// EnterAlterView is called when entering the alterView production.
	EnterAlterView(c *AlterViewContext)

	// EnterAlterStorageVault is called when entering the alterStorageVault production.
	EnterAlterStorageVault(c *AlterStorageVaultContext)

	// EnterDropCatalogRecycleBin is called when entering the dropCatalogRecycleBin production.
	EnterDropCatalogRecycleBin(c *DropCatalogRecycleBinContext)

	// EnterAdminShowReplicaStatus is called when entering the adminShowReplicaStatus production.
	EnterAdminShowReplicaStatus(c *AdminShowReplicaStatusContext)

	// EnterAdminShowReplicaDistribution is called when entering the adminShowReplicaDistribution production.
	EnterAdminShowReplicaDistribution(c *AdminShowReplicaDistributionContext)

	// EnterAdminSetReplicaStatus is called when entering the adminSetReplicaStatus production.
	EnterAdminSetReplicaStatus(c *AdminSetReplicaStatusContext)

	// EnterAdminSetReplicaVersion is called when entering the adminSetReplicaVersion production.
	EnterAdminSetReplicaVersion(c *AdminSetReplicaVersionContext)

	// EnterAdminRepairTable is called when entering the adminRepairTable production.
	EnterAdminRepairTable(c *AdminRepairTableContext)

	// EnterAdminCancelRepairTable is called when entering the adminCancelRepairTable production.
	EnterAdminCancelRepairTable(c *AdminCancelRepairTableContext)

	// EnterAdminCompactTable is called when entering the adminCompactTable production.
	EnterAdminCompactTable(c *AdminCompactTableContext)

	// EnterAdminSetFrontendConfig is called when entering the adminSetFrontendConfig production.
	EnterAdminSetFrontendConfig(c *AdminSetFrontendConfigContext)

	// EnterAdminCheckTablets is called when entering the adminCheckTablets production.
	EnterAdminCheckTablets(c *AdminCheckTabletsContext)

	// EnterAdminRebalanceDisk is called when entering the adminRebalanceDisk production.
	EnterAdminRebalanceDisk(c *AdminRebalanceDiskContext)

	// EnterAdminCancelRebalanceDisk is called when entering the adminCancelRebalanceDisk production.
	EnterAdminCancelRebalanceDisk(c *AdminCancelRebalanceDiskContext)

	// EnterAdminCleanTrash is called when entering the adminCleanTrash production.
	EnterAdminCleanTrash(c *AdminCleanTrashContext)

	// EnterAdminSetPartitionVersion is called when entering the adminSetPartitionVersion production.
	EnterAdminSetPartitionVersion(c *AdminSetPartitionVersionContext)

	// EnterAdminDiagnoseTablet is called when entering the adminDiagnoseTablet production.
	EnterAdminDiagnoseTablet(c *AdminDiagnoseTabletContext)

	// EnterAdminShowTabletStorageFormat is called when entering the adminShowTabletStorageFormat production.
	EnterAdminShowTabletStorageFormat(c *AdminShowTabletStorageFormatContext)

	// EnterAdminCopyTablet is called when entering the adminCopyTablet production.
	EnterAdminCopyTablet(c *AdminCopyTabletContext)

	// EnterAdminSetTableStatus is called when entering the adminSetTableStatus production.
	EnterAdminSetTableStatus(c *AdminSetTableStatusContext)

	// EnterBaseTableRef is called when entering the baseTableRef production.
	EnterBaseTableRef(c *BaseTableRefContext)

	// EnterWildWhere is called when entering the wildWhere production.
	EnterWildWhere(c *WildWhereContext)

	// EnterTransactionBegin is called when entering the transactionBegin production.
	EnterTransactionBegin(c *TransactionBeginContext)

	// EnterTranscationCommit is called when entering the transcationCommit production.
	EnterTranscationCommit(c *TranscationCommitContext)

	// EnterTransactionRollback is called when entering the transactionRollback production.
	EnterTransactionRollback(c *TransactionRollbackContext)

	// EnterGrantTablePrivilege is called when entering the grantTablePrivilege production.
	EnterGrantTablePrivilege(c *GrantTablePrivilegeContext)

	// EnterGrantResourcePrivilege is called when entering the grantResourcePrivilege production.
	EnterGrantResourcePrivilege(c *GrantResourcePrivilegeContext)

	// EnterGrantRole is called when entering the grantRole production.
	EnterGrantRole(c *GrantRoleContext)

	// EnterPrivilege is called when entering the privilege production.
	EnterPrivilege(c *PrivilegeContext)

	// EnterPrivilegeList is called when entering the privilegeList production.
	EnterPrivilegeList(c *PrivilegeListContext)

	// EnterAlterTable is called when entering the alterTable production.
	EnterAlterTable(c *AlterTableContext)

	// EnterAlterTableAddRollup is called when entering the alterTableAddRollup production.
	EnterAlterTableAddRollup(c *AlterTableAddRollupContext)

	// EnterAlterTableDropRollup is called when entering the alterTableDropRollup production.
	EnterAlterTableDropRollup(c *AlterTableDropRollupContext)

	// EnterAlterSystem is called when entering the alterSystem production.
	EnterAlterSystem(c *AlterSystemContext)

	// EnterAlterDatabaseSetQuota is called when entering the alterDatabaseSetQuota production.
	EnterAlterDatabaseSetQuota(c *AlterDatabaseSetQuotaContext)

	// EnterAlterDatabaseRename is called when entering the alterDatabaseRename production.
	EnterAlterDatabaseRename(c *AlterDatabaseRenameContext)

	// EnterAlterDatabaseProperties is called when entering the alterDatabaseProperties production.
	EnterAlterDatabaseProperties(c *AlterDatabasePropertiesContext)

	// EnterAlterCatalogRename is called when entering the alterCatalogRename production.
	EnterAlterCatalogRename(c *AlterCatalogRenameContext)

	// EnterAlterCatalogProperties is called when entering the alterCatalogProperties production.
	EnterAlterCatalogProperties(c *AlterCatalogPropertiesContext)

	// EnterAlterCatalogComment is called when entering the alterCatalogComment production.
	EnterAlterCatalogComment(c *AlterCatalogCommentContext)

	// EnterAlterResource is called when entering the alterResource production.
	EnterAlterResource(c *AlterResourceContext)

	// EnterAlterColocateGroup is called when entering the alterColocateGroup production.
	EnterAlterColocateGroup(c *AlterColocateGroupContext)

	// EnterAlterWorkloadGroup is called when entering the alterWorkloadGroup production.
	EnterAlterWorkloadGroup(c *AlterWorkloadGroupContext)

	// EnterAlterWorkloadPolicy is called when entering the alterWorkloadPolicy production.
	EnterAlterWorkloadPolicy(c *AlterWorkloadPolicyContext)

	// EnterAlterRoutineLoad is called when entering the alterRoutineLoad production.
	EnterAlterRoutineLoad(c *AlterRoutineLoadContext)

	// EnterAlterSqlBlockRule is called when entering the alterSqlBlockRule production.
	EnterAlterSqlBlockRule(c *AlterSqlBlockRuleContext)

	// EnterAlterTableProperties is called when entering the alterTableProperties production.
	EnterAlterTableProperties(c *AlterTablePropertiesContext)

	// EnterAlterStoragePlicy is called when entering the alterStoragePlicy production.
	EnterAlterStoragePlicy(c *AlterStoragePlicyContext)

	// EnterAlterUser is called when entering the alterUser production.
	EnterAlterUser(c *AlterUserContext)

	// EnterAlterRepository is called when entering the alterRepository production.
	EnterAlterRepository(c *AlterRepositoryContext)

	// EnterAddBackendClause is called when entering the addBackendClause production.
	EnterAddBackendClause(c *AddBackendClauseContext)

	// EnterDropBackendClause is called when entering the dropBackendClause production.
	EnterDropBackendClause(c *DropBackendClauseContext)

	// EnterDecommissionBackendClause is called when entering the decommissionBackendClause production.
	EnterDecommissionBackendClause(c *DecommissionBackendClauseContext)

	// EnterAddObserverClause is called when entering the addObserverClause production.
	EnterAddObserverClause(c *AddObserverClauseContext)

	// EnterDropObserverClause is called when entering the dropObserverClause production.
	EnterDropObserverClause(c *DropObserverClauseContext)

	// EnterAddFollowerClause is called when entering the addFollowerClause production.
	EnterAddFollowerClause(c *AddFollowerClauseContext)

	// EnterDropFollowerClause is called when entering the dropFollowerClause production.
	EnterDropFollowerClause(c *DropFollowerClauseContext)

	// EnterAddBrokerClause is called when entering the addBrokerClause production.
	EnterAddBrokerClause(c *AddBrokerClauseContext)

	// EnterDropBrokerClause is called when entering the dropBrokerClause production.
	EnterDropBrokerClause(c *DropBrokerClauseContext)

	// EnterDropAllBrokerClause is called when entering the dropAllBrokerClause production.
	EnterDropAllBrokerClause(c *DropAllBrokerClauseContext)

	// EnterAlterLoadErrorUrlClause is called when entering the alterLoadErrorUrlClause production.
	EnterAlterLoadErrorUrlClause(c *AlterLoadErrorUrlClauseContext)

	// EnterModifyBackendClause is called when entering the modifyBackendClause production.
	EnterModifyBackendClause(c *ModifyBackendClauseContext)

	// EnterModifyFrontendOrBackendHostNameClause is called when entering the modifyFrontendOrBackendHostNameClause production.
	EnterModifyFrontendOrBackendHostNameClause(c *ModifyFrontendOrBackendHostNameClauseContext)

	// EnterDropRollupClause is called when entering the dropRollupClause production.
	EnterDropRollupClause(c *DropRollupClauseContext)

	// EnterAddRollupClause is called when entering the addRollupClause production.
	EnterAddRollupClause(c *AddRollupClauseContext)

	// EnterAddColumnClause is called when entering the addColumnClause production.
	EnterAddColumnClause(c *AddColumnClauseContext)

	// EnterAddColumnsClause is called when entering the addColumnsClause production.
	EnterAddColumnsClause(c *AddColumnsClauseContext)

	// EnterDropColumnClause is called when entering the dropColumnClause production.
	EnterDropColumnClause(c *DropColumnClauseContext)

	// EnterModifyColumnClause is called when entering the modifyColumnClause production.
	EnterModifyColumnClause(c *ModifyColumnClauseContext)

	// EnterReorderColumnsClause is called when entering the reorderColumnsClause production.
	EnterReorderColumnsClause(c *ReorderColumnsClauseContext)

	// EnterAddPartitionClause is called when entering the addPartitionClause production.
	EnterAddPartitionClause(c *AddPartitionClauseContext)

	// EnterDropPartitionClause is called when entering the dropPartitionClause production.
	EnterDropPartitionClause(c *DropPartitionClauseContext)

	// EnterModifyPartitionClause is called when entering the modifyPartitionClause production.
	EnterModifyPartitionClause(c *ModifyPartitionClauseContext)

	// EnterReplacePartitionClause is called when entering the replacePartitionClause production.
	EnterReplacePartitionClause(c *ReplacePartitionClauseContext)

	// EnterReplaceTableClause is called when entering the replaceTableClause production.
	EnterReplaceTableClause(c *ReplaceTableClauseContext)

	// EnterRenameClause is called when entering the renameClause production.
	EnterRenameClause(c *RenameClauseContext)

	// EnterRenameRollupClause is called when entering the renameRollupClause production.
	EnterRenameRollupClause(c *RenameRollupClauseContext)

	// EnterRenamePartitionClause is called when entering the renamePartitionClause production.
	EnterRenamePartitionClause(c *RenamePartitionClauseContext)

	// EnterRenameColumnClause is called when entering the renameColumnClause production.
	EnterRenameColumnClause(c *RenameColumnClauseContext)

	// EnterAddIndexClause is called when entering the addIndexClause production.
	EnterAddIndexClause(c *AddIndexClauseContext)

	// EnterDropIndexClause is called when entering the dropIndexClause production.
	EnterDropIndexClause(c *DropIndexClauseContext)

	// EnterEnableFeatureClause is called when entering the enableFeatureClause production.
	EnterEnableFeatureClause(c *EnableFeatureClauseContext)

	// EnterModifyDistributionClause is called when entering the modifyDistributionClause production.
	EnterModifyDistributionClause(c *ModifyDistributionClauseContext)

	// EnterModifyTableCommentClause is called when entering the modifyTableCommentClause production.
	EnterModifyTableCommentClause(c *ModifyTableCommentClauseContext)

	// EnterModifyColumnCommentClause is called when entering the modifyColumnCommentClause production.
	EnterModifyColumnCommentClause(c *ModifyColumnCommentClauseContext)

	// EnterModifyEngineClause is called when entering the modifyEngineClause production.
	EnterModifyEngineClause(c *ModifyEngineClauseContext)

	// EnterAlterMultiPartitionClause is called when entering the alterMultiPartitionClause production.
	EnterAlterMultiPartitionClause(c *AlterMultiPartitionClauseContext)

	// EnterColumnPosition is called when entering the columnPosition production.
	EnterColumnPosition(c *ColumnPositionContext)

	// EnterToRollup is called when entering the toRollup production.
	EnterToRollup(c *ToRollupContext)

	// EnterFromRollup is called when entering the fromRollup production.
	EnterFromRollup(c *FromRollupContext)

	// EnterDropDatabase is called when entering the dropDatabase production.
	EnterDropDatabase(c *DropDatabaseContext)

	// EnterDropCatalog is called when entering the dropCatalog production.
	EnterDropCatalog(c *DropCatalogContext)

	// EnterDropFunction is called when entering the dropFunction production.
	EnterDropFunction(c *DropFunctionContext)

	// EnterDropTable is called when entering the dropTable production.
	EnterDropTable(c *DropTableContext)

	// EnterDropUser is called when entering the dropUser production.
	EnterDropUser(c *DropUserContext)

	// EnterDropView is called when entering the dropView production.
	EnterDropView(c *DropViewContext)

	// EnterDropRepository is called when entering the dropRepository production.
	EnterDropRepository(c *DropRepositoryContext)

	// EnterDropRole is called when entering the dropRole production.
	EnterDropRole(c *DropRoleContext)

	// EnterDropFile is called when entering the dropFile production.
	EnterDropFile(c *DropFileContext)

	// EnterDropIndex is called when entering the dropIndex production.
	EnterDropIndex(c *DropIndexContext)

	// EnterDropResource is called when entering the dropResource production.
	EnterDropResource(c *DropResourceContext)

	// EnterDropWorkloadGroup is called when entering the dropWorkloadGroup production.
	EnterDropWorkloadGroup(c *DropWorkloadGroupContext)

	// EnterDropWorkloadPolicy is called when entering the dropWorkloadPolicy production.
	EnterDropWorkloadPolicy(c *DropWorkloadPolicyContext)

	// EnterDropEncryptkey is called when entering the dropEncryptkey production.
	EnterDropEncryptkey(c *DropEncryptkeyContext)

	// EnterDropSqlBlockRule is called when entering the dropSqlBlockRule production.
	EnterDropSqlBlockRule(c *DropSqlBlockRuleContext)

	// EnterDropRowPolicy is called when entering the dropRowPolicy production.
	EnterDropRowPolicy(c *DropRowPolicyContext)

	// EnterDropStoragePolicy is called when entering the dropStoragePolicy production.
	EnterDropStoragePolicy(c *DropStoragePolicyContext)

	// EnterDropStage is called when entering the dropStage production.
	EnterDropStage(c *DropStageContext)

	// EnterAlterTableStats is called when entering the alterTableStats production.
	EnterAlterTableStats(c *AlterTableStatsContext)

	// EnterAlterColumnStats is called when entering the alterColumnStats production.
	EnterAlterColumnStats(c *AlterColumnStatsContext)

	// EnterDropStats is called when entering the dropStats production.
	EnterDropStats(c *DropStatsContext)

	// EnterDropCachedStats is called when entering the dropCachedStats production.
	EnterDropCachedStats(c *DropCachedStatsContext)

	// EnterDropExpiredStats is called when entering the dropExpiredStats production.
	EnterDropExpiredStats(c *DropExpiredStatsContext)

	// EnterDropAanalyzeJob is called when entering the dropAanalyzeJob production.
	EnterDropAanalyzeJob(c *DropAanalyzeJobContext)

	// EnterCreateDatabase is called when entering the createDatabase production.
	EnterCreateDatabase(c *CreateDatabaseContext)

	// EnterCreateCatalog is called when entering the createCatalog production.
	EnterCreateCatalog(c *CreateCatalogContext)

	// EnterCreateUserDefineFunction is called when entering the createUserDefineFunction production.
	EnterCreateUserDefineFunction(c *CreateUserDefineFunctionContext)

	// EnterCreateAliasFunction is called when entering the createAliasFunction production.
	EnterCreateAliasFunction(c *CreateAliasFunctionContext)

	// EnterCreateUser is called when entering the createUser production.
	EnterCreateUser(c *CreateUserContext)

	// EnterCreateRepository is called when entering the createRepository production.
	EnterCreateRepository(c *CreateRepositoryContext)

	// EnterCreateRole is called when entering the createRole production.
	EnterCreateRole(c *CreateRoleContext)

	// EnterCreateFile is called when entering the createFile production.
	EnterCreateFile(c *CreateFileContext)

	// EnterCreateIndex is called when entering the createIndex production.
	EnterCreateIndex(c *CreateIndexContext)

	// EnterCreateResource is called when entering the createResource production.
	EnterCreateResource(c *CreateResourceContext)

	// EnterCreateStorageVault is called when entering the createStorageVault production.
	EnterCreateStorageVault(c *CreateStorageVaultContext)

	// EnterCreateWorkloadGroup is called when entering the createWorkloadGroup production.
	EnterCreateWorkloadGroup(c *CreateWorkloadGroupContext)

	// EnterCreateWorkloadPolicy is called when entering the createWorkloadPolicy production.
	EnterCreateWorkloadPolicy(c *CreateWorkloadPolicyContext)

	// EnterCreateEncryptkey is called when entering the createEncryptkey production.
	EnterCreateEncryptkey(c *CreateEncryptkeyContext)

	// EnterCreateDataSyncJob is called when entering the createDataSyncJob production.
	EnterCreateDataSyncJob(c *CreateDataSyncJobContext)

	// EnterCreateSqlBlockRule is called when entering the createSqlBlockRule production.
	EnterCreateSqlBlockRule(c *CreateSqlBlockRuleContext)

	// EnterCreateStoragePolicy is called when entering the createStoragePolicy production.
	EnterCreateStoragePolicy(c *CreateStoragePolicyContext)

	// EnterBuildIndex is called when entering the buildIndex production.
	EnterBuildIndex(c *BuildIndexContext)

	// EnterCreateStage is called when entering the createStage production.
	EnterCreateStage(c *CreateStageContext)

	// EnterChannelDescriptions is called when entering the channelDescriptions production.
	EnterChannelDescriptions(c *ChannelDescriptionsContext)

	// EnterChannelDescription is called when entering the channelDescription production.
	EnterChannelDescription(c *ChannelDescriptionContext)

	// EnterWorkloadPolicyActions is called when entering the workloadPolicyActions production.
	EnterWorkloadPolicyActions(c *WorkloadPolicyActionsContext)

	// EnterWorkloadPolicyAction is called when entering the workloadPolicyAction production.
	EnterWorkloadPolicyAction(c *WorkloadPolicyActionContext)

	// EnterWorkloadPolicyConditions is called when entering the workloadPolicyConditions production.
	EnterWorkloadPolicyConditions(c *WorkloadPolicyConditionsContext)

	// EnterWorkloadPolicyCondition is called when entering the workloadPolicyCondition production.
	EnterWorkloadPolicyCondition(c *WorkloadPolicyConditionContext)

	// EnterStorageBackend is called when entering the storageBackend production.
	EnterStorageBackend(c *StorageBackendContext)

	// EnterPasswordOption is called when entering the passwordOption production.
	EnterPasswordOption(c *PasswordOptionContext)

	// EnterFunctionArguments is called when entering the functionArguments production.
	EnterFunctionArguments(c *FunctionArgumentsContext)

	// EnterFunctionArgument is called when entering the functionArgument production.
	EnterFunctionArgument(c *FunctionArgumentContext)

	// EnterSetOptions is called when entering the setOptions production.
	EnterSetOptions(c *SetOptionsContext)

	// EnterSetDefaultStorageVault is called when entering the setDefaultStorageVault production.
	EnterSetDefaultStorageVault(c *SetDefaultStorageVaultContext)

	// EnterSetUserProperties is called when entering the setUserProperties production.
	EnterSetUserProperties(c *SetUserPropertiesContext)

	// EnterSetTransaction is called when entering the setTransaction production.
	EnterSetTransaction(c *SetTransactionContext)

	// EnterOptionWithType is called when entering the optionWithType production.
	EnterOptionWithType(c *OptionWithTypeContext)

	// EnterSetNames is called when entering the setNames production.
	EnterSetNames(c *SetNamesContext)

	// EnterSetCharset is called when entering the setCharset production.
	EnterSetCharset(c *SetCharsetContext)

	// EnterSetCollate is called when entering the setCollate production.
	EnterSetCollate(c *SetCollateContext)

	// EnterSetPassword is called when entering the setPassword production.
	EnterSetPassword(c *SetPasswordContext)

	// EnterSetLdapAdminPassword is called when entering the setLdapAdminPassword production.
	EnterSetLdapAdminPassword(c *SetLdapAdminPasswordContext)

	// EnterSetVariableWithoutType is called when entering the setVariableWithoutType production.
	EnterSetVariableWithoutType(c *SetVariableWithoutTypeContext)

	// EnterSetSystemVariable is called when entering the setSystemVariable production.
	EnterSetSystemVariable(c *SetSystemVariableContext)

	// EnterSetUserVariable is called when entering the setUserVariable production.
	EnterSetUserVariable(c *SetUserVariableContext)

	// EnterTransactionAccessMode is called when entering the transactionAccessMode production.
	EnterTransactionAccessMode(c *TransactionAccessModeContext)

	// EnterIsolationLevel is called when entering the isolationLevel production.
	EnterIsolationLevel(c *IsolationLevelContext)

	// EnterUnsupoortedUnsetStatement is called when entering the unsupoortedUnsetStatement production.
	EnterUnsupoortedUnsetStatement(c *UnsupoortedUnsetStatementContext)

	// EnterUseDatabase is called when entering the useDatabase production.
	EnterUseDatabase(c *UseDatabaseContext)

	// EnterUseCloudCluster is called when entering the useCloudCluster production.
	EnterUseCloudCluster(c *UseCloudClusterContext)

	// EnterSwitchCatalog is called when entering the switchCatalog production.
	EnterSwitchCatalog(c *SwitchCatalogContext)

	// EnterTruncateTable is called when entering the truncateTable production.
	EnterTruncateTable(c *TruncateTableContext)

	// EnterKillConnection is called when entering the killConnection production.
	EnterKillConnection(c *KillConnectionContext)

	// EnterKillQuery is called when entering the killQuery production.
	EnterKillQuery(c *KillQueryContext)

	// EnterDescribeTableValuedFunction is called when entering the describeTableValuedFunction production.
	EnterDescribeTableValuedFunction(c *DescribeTableValuedFunctionContext)

	// EnterDescribeTableAll is called when entering the describeTableAll production.
	EnterDescribeTableAll(c *DescribeTableAllContext)

	// EnterDescribeTable is called when entering the describeTable production.
	EnterDescribeTable(c *DescribeTableContext)

	// EnterConstraint is called when entering the constraint production.
	EnterConstraint(c *ConstraintContext)

	// EnterPartitionSpec is called when entering the partitionSpec production.
	EnterPartitionSpec(c *PartitionSpecContext)

	// EnterPartitionTable is called when entering the partitionTable production.
	EnterPartitionTable(c *PartitionTableContext)

	// EnterIdentityOrFunctionList is called when entering the identityOrFunctionList production.
	EnterIdentityOrFunctionList(c *IdentityOrFunctionListContext)

	// EnterIdentityOrFunction is called when entering the identityOrFunction production.
	EnterIdentityOrFunction(c *IdentityOrFunctionContext)

	// EnterDataDesc is called when entering the dataDesc production.
	EnterDataDesc(c *DataDescContext)

	// EnterBuildMode is called when entering the buildMode production.
	EnterBuildMode(c *BuildModeContext)

	// EnterRefreshTrigger is called when entering the refreshTrigger production.
	EnterRefreshTrigger(c *RefreshTriggerContext)

	// EnterRefreshSchedule is called when entering the refreshSchedule production.
	EnterRefreshSchedule(c *RefreshScheduleContext)

	// EnterRefreshMethod is called when entering the refreshMethod production.
	EnterRefreshMethod(c *RefreshMethodContext)

	// EnterMvPartition is called when entering the mvPartition production.
	EnterMvPartition(c *MvPartitionContext)

	// EnterIdentifierOrText is called when entering the identifierOrText production.
	EnterIdentifierOrText(c *IdentifierOrTextContext)

	// EnterIdentifierOrTextOrAsterisk is called when entering the identifierOrTextOrAsterisk production.
	EnterIdentifierOrTextOrAsterisk(c *IdentifierOrTextOrAsteriskContext)

	// EnterMultipartIdentifierOrAsterisk is called when entering the multipartIdentifierOrAsterisk production.
	EnterMultipartIdentifierOrAsterisk(c *MultipartIdentifierOrAsteriskContext)

	// EnterIdentifierOrAsterisk is called when entering the identifierOrAsterisk production.
	EnterIdentifierOrAsterisk(c *IdentifierOrAsteriskContext)

	// EnterUserIdentify is called when entering the userIdentify production.
	EnterUserIdentify(c *UserIdentifyContext)

	// EnterGrantUserIdentify is called when entering the grantUserIdentify production.
	EnterGrantUserIdentify(c *GrantUserIdentifyContext)

	// EnterExplain is called when entering the explain production.
	EnterExplain(c *ExplainContext)

	// EnterExplainCommand is called when entering the explainCommand production.
	EnterExplainCommand(c *ExplainCommandContext)

	// EnterPlanType is called when entering the planType production.
	EnterPlanType(c *PlanTypeContext)

	// EnterMergeType is called when entering the mergeType production.
	EnterMergeType(c *MergeTypeContext)

	// EnterPreFilterClause is called when entering the preFilterClause production.
	EnterPreFilterClause(c *PreFilterClauseContext)

	// EnterDeleteOnClause is called when entering the deleteOnClause production.
	EnterDeleteOnClause(c *DeleteOnClauseContext)

	// EnterSequenceColClause is called when entering the sequenceColClause production.
	EnterSequenceColClause(c *SequenceColClauseContext)

	// EnterColFromPath is called when entering the colFromPath production.
	EnterColFromPath(c *ColFromPathContext)

	// EnterColMappingList is called when entering the colMappingList production.
	EnterColMappingList(c *ColMappingListContext)

	// EnterMappingExpr is called when entering the mappingExpr production.
	EnterMappingExpr(c *MappingExprContext)

	// EnterWithRemoteStorageSystem is called when entering the withRemoteStorageSystem production.
	EnterWithRemoteStorageSystem(c *WithRemoteStorageSystemContext)

	// EnterResourceDesc is called when entering the resourceDesc production.
	EnterResourceDesc(c *ResourceDescContext)

	// EnterMysqlDataDesc is called when entering the mysqlDataDesc production.
	EnterMysqlDataDesc(c *MysqlDataDescContext)

	// EnterSkipLines is called when entering the skipLines production.
	EnterSkipLines(c *SkipLinesContext)

	// EnterOutFileClause is called when entering the outFileClause production.
	EnterOutFileClause(c *OutFileClauseContext)

	// EnterQuery is called when entering the query production.
	EnterQuery(c *QueryContext)

	// EnterQueryTermDefault is called when entering the queryTermDefault production.
	EnterQueryTermDefault(c *QueryTermDefaultContext)

	// EnterSetOperation is called when entering the setOperation production.
	EnterSetOperation(c *SetOperationContext)

	// EnterSetQuantifier is called when entering the setQuantifier production.
	EnterSetQuantifier(c *SetQuantifierContext)

	// EnterQueryPrimaryDefault is called when entering the queryPrimaryDefault production.
	EnterQueryPrimaryDefault(c *QueryPrimaryDefaultContext)

	// EnterSubquery is called when entering the subquery production.
	EnterSubquery(c *SubqueryContext)

	// EnterValuesTable is called when entering the valuesTable production.
	EnterValuesTable(c *ValuesTableContext)

	// EnterRegularQuerySpecification is called when entering the regularQuerySpecification production.
	EnterRegularQuerySpecification(c *RegularQuerySpecificationContext)

	// EnterCte is called when entering the cte production.
	EnterCte(c *CteContext)

	// EnterAliasQuery is called when entering the aliasQuery production.
	EnterAliasQuery(c *AliasQueryContext)

	// EnterColumnAliases is called when entering the columnAliases production.
	EnterColumnAliases(c *ColumnAliasesContext)

	// EnterSelectClause is called when entering the selectClause production.
	EnterSelectClause(c *SelectClauseContext)

	// EnterSelectColumnClause is called when entering the selectColumnClause production.
	EnterSelectColumnClause(c *SelectColumnClauseContext)

	// EnterWhereClause is called when entering the whereClause production.
	EnterWhereClause(c *WhereClauseContext)

	// EnterFromClause is called when entering the fromClause production.
	EnterFromClause(c *FromClauseContext)

	// EnterIntoClause is called when entering the intoClause production.
	EnterIntoClause(c *IntoClauseContext)

	// EnterBulkCollectClause is called when entering the bulkCollectClause production.
	EnterBulkCollectClause(c *BulkCollectClauseContext)

	// EnterTableRow is called when entering the tableRow production.
	EnterTableRow(c *TableRowContext)

	// EnterRelations is called when entering the relations production.
	EnterRelations(c *RelationsContext)

	// EnterRelation is called when entering the relation production.
	EnterRelation(c *RelationContext)

	// EnterJoinRelation is called when entering the joinRelation production.
	EnterJoinRelation(c *JoinRelationContext)

	// EnterBracketDistributeType is called when entering the bracketDistributeType production.
	EnterBracketDistributeType(c *BracketDistributeTypeContext)

	// EnterCommentDistributeType is called when entering the commentDistributeType production.
	EnterCommentDistributeType(c *CommentDistributeTypeContext)

	// EnterBracketRelationHint is called when entering the bracketRelationHint production.
	EnterBracketRelationHint(c *BracketRelationHintContext)

	// EnterCommentRelationHint is called when entering the commentRelationHint production.
	EnterCommentRelationHint(c *CommentRelationHintContext)

	// EnterAggClause is called when entering the aggClause production.
	EnterAggClause(c *AggClauseContext)

	// EnterGroupingElement is called when entering the groupingElement production.
	EnterGroupingElement(c *GroupingElementContext)

	// EnterGroupingSet is called when entering the groupingSet production.
	EnterGroupingSet(c *GroupingSetContext)

	// EnterHavingClause is called when entering the havingClause production.
	EnterHavingClause(c *HavingClauseContext)

	// EnterSelectHint is called when entering the selectHint production.
	EnterSelectHint(c *SelectHintContext)

	// EnterHintStatement is called when entering the hintStatement production.
	EnterHintStatement(c *HintStatementContext)

	// EnterHintAssignment is called when entering the hintAssignment production.
	EnterHintAssignment(c *HintAssignmentContext)

	// EnterUpdateAssignment is called when entering the updateAssignment production.
	EnterUpdateAssignment(c *UpdateAssignmentContext)

	// EnterUpdateAssignmentSeq is called when entering the updateAssignmentSeq production.
	EnterUpdateAssignmentSeq(c *UpdateAssignmentSeqContext)

	// EnterLateralView is called when entering the lateralView production.
	EnterLateralView(c *LateralViewContext)

	// EnterQueryOrganization is called when entering the queryOrganization production.
	EnterQueryOrganization(c *QueryOrganizationContext)

	// EnterSortClause is called when entering the sortClause production.
	EnterSortClause(c *SortClauseContext)

	// EnterSortItem is called when entering the sortItem production.
	EnterSortItem(c *SortItemContext)

	// EnterLimitClause is called when entering the limitClause production.
	EnterLimitClause(c *LimitClauseContext)

	// EnterPartitionClause is called when entering the partitionClause production.
	EnterPartitionClause(c *PartitionClauseContext)

	// EnterJoinType is called when entering the joinType production.
	EnterJoinType(c *JoinTypeContext)

	// EnterJoinCriteria is called when entering the joinCriteria production.
	EnterJoinCriteria(c *JoinCriteriaContext)

	// EnterIdentifierList is called when entering the identifierList production.
	EnterIdentifierList(c *IdentifierListContext)

	// EnterIdentifierSeq is called when entering the identifierSeq production.
	EnterIdentifierSeq(c *IdentifierSeqContext)

	// EnterOptScanParams is called when entering the optScanParams production.
	EnterOptScanParams(c *OptScanParamsContext)

	// EnterTableName is called when entering the tableName production.
	EnterTableName(c *TableNameContext)

	// EnterAliasedQuery is called when entering the aliasedQuery production.
	EnterAliasedQuery(c *AliasedQueryContext)

	// EnterTableValuedFunction is called when entering the tableValuedFunction production.
	EnterTableValuedFunction(c *TableValuedFunctionContext)

	// EnterRelationList is called when entering the relationList production.
	EnterRelationList(c *RelationListContext)

	// EnterMaterializedViewName is called when entering the materializedViewName production.
	EnterMaterializedViewName(c *MaterializedViewNameContext)

	// EnterPropertyClause is called when entering the propertyClause production.
	EnterPropertyClause(c *PropertyClauseContext)

	// EnterPropertyItemList is called when entering the propertyItemList production.
	EnterPropertyItemList(c *PropertyItemListContext)

	// EnterPropertyItem is called when entering the propertyItem production.
	EnterPropertyItem(c *PropertyItemContext)

	// EnterPropertyKey is called when entering the propertyKey production.
	EnterPropertyKey(c *PropertyKeyContext)

	// EnterPropertyValue is called when entering the propertyValue production.
	EnterPropertyValue(c *PropertyValueContext)

	// EnterTableAlias is called when entering the tableAlias production.
	EnterTableAlias(c *TableAliasContext)

	// EnterMultipartIdentifier is called when entering the multipartIdentifier production.
	EnterMultipartIdentifier(c *MultipartIdentifierContext)

	// EnterSimpleColumnDefs is called when entering the simpleColumnDefs production.
	EnterSimpleColumnDefs(c *SimpleColumnDefsContext)

	// EnterSimpleColumnDef is called when entering the simpleColumnDef production.
	EnterSimpleColumnDef(c *SimpleColumnDefContext)

	// EnterColumnDefs is called when entering the columnDefs production.
	EnterColumnDefs(c *ColumnDefsContext)

	// EnterColumnDef is called when entering the columnDef production.
	EnterColumnDef(c *ColumnDefContext)

	// EnterIndexDefs is called when entering the indexDefs production.
	EnterIndexDefs(c *IndexDefsContext)

	// EnterIndexDef is called when entering the indexDef production.
	EnterIndexDef(c *IndexDefContext)

	// EnterPartitionsDef is called when entering the partitionsDef production.
	EnterPartitionsDef(c *PartitionsDefContext)

	// EnterPartitionDef is called when entering the partitionDef production.
	EnterPartitionDef(c *PartitionDefContext)

	// EnterLessThanPartitionDef is called when entering the lessThanPartitionDef production.
	EnterLessThanPartitionDef(c *LessThanPartitionDefContext)

	// EnterFixedPartitionDef is called when entering the fixedPartitionDef production.
	EnterFixedPartitionDef(c *FixedPartitionDefContext)

	// EnterStepPartitionDef is called when entering the stepPartitionDef production.
	EnterStepPartitionDef(c *StepPartitionDefContext)

	// EnterInPartitionDef is called when entering the inPartitionDef production.
	EnterInPartitionDef(c *InPartitionDefContext)

	// EnterPartitionValueList is called when entering the partitionValueList production.
	EnterPartitionValueList(c *PartitionValueListContext)

	// EnterPartitionValueDef is called when entering the partitionValueDef production.
	EnterPartitionValueDef(c *PartitionValueDefContext)

	// EnterRollupDefs is called when entering the rollupDefs production.
	EnterRollupDefs(c *RollupDefsContext)

	// EnterRollupDef is called when entering the rollupDef production.
	EnterRollupDef(c *RollupDefContext)

	// EnterAggTypeDef is called when entering the aggTypeDef production.
	EnterAggTypeDef(c *AggTypeDefContext)

	// EnterTabletList is called when entering the tabletList production.
	EnterTabletList(c *TabletListContext)

	// EnterInlineTable is called when entering the inlineTable production.
	EnterInlineTable(c *InlineTableContext)

	// EnterNamedExpression is called when entering the namedExpression production.
	EnterNamedExpression(c *NamedExpressionContext)

	// EnterNamedExpressionSeq is called when entering the namedExpressionSeq production.
	EnterNamedExpressionSeq(c *NamedExpressionSeqContext)

	// EnterExpression is called when entering the expression production.
	EnterExpression(c *ExpressionContext)

	// EnterLambdaExpression is called when entering the lambdaExpression production.
	EnterLambdaExpression(c *LambdaExpressionContext)

	// EnterExist is called when entering the exist production.
	EnterExist(c *ExistContext)

	// EnterLogicalNot is called when entering the logicalNot production.
	EnterLogicalNot(c *LogicalNotContext)

	// EnterPredicated is called when entering the predicated production.
	EnterPredicated(c *PredicatedContext)

	// EnterIsnull is called when entering the isnull production.
	EnterIsnull(c *IsnullContext)

	// EnterIs_not_null_pred is called when entering the is_not_null_pred production.
	EnterIs_not_null_pred(c *Is_not_null_predContext)

	// EnterLogicalBinary is called when entering the logicalBinary production.
	EnterLogicalBinary(c *LogicalBinaryContext)

	// EnterDoublePipes is called when entering the doublePipes production.
	EnterDoublePipes(c *DoublePipesContext)

	// EnterRowConstructor is called when entering the rowConstructor production.
	EnterRowConstructor(c *RowConstructorContext)

	// EnterRowConstructorItem is called when entering the rowConstructorItem production.
	EnterRowConstructorItem(c *RowConstructorItemContext)

	// EnterPredicate is called when entering the predicate production.
	EnterPredicate(c *PredicateContext)

	// EnterValueExpressionDefault is called when entering the valueExpressionDefault production.
	EnterValueExpressionDefault(c *ValueExpressionDefaultContext)

	// EnterComparison is called when entering the comparison production.
	EnterComparison(c *ComparisonContext)

	// EnterBitOperation is called when entering the bitOperation production.
	EnterBitOperation(c *BitOperationContext)

	// EnterArithmeticBinary is called when entering the arithmeticBinary production.
	EnterArithmeticBinary(c *ArithmeticBinaryContext)

	// EnterArithmeticUnary is called when entering the arithmeticUnary production.
	EnterArithmeticUnary(c *ArithmeticUnaryContext)

	// EnterDatetimeUnit is called when entering the datetimeUnit production.
	EnterDatetimeUnit(c *DatetimeUnitContext)

	// EnterDateCeil is called when entering the dateCeil production.
	EnterDateCeil(c *DateCeilContext)

	// EnterDereference is called when entering the dereference production.
	EnterDereference(c *DereferenceContext)

	// EnterCurrentDate is called when entering the currentDate production.
	EnterCurrentDate(c *CurrentDateContext)

	// EnterTimestampadd is called when entering the timestampadd production.
	EnterTimestampadd(c *TimestampaddContext)

	// EnterDate_sub is called when entering the date_sub production.
	EnterDate_sub(c *Date_subContext)

	// EnterCast is called when entering the cast production.
	EnterCast(c *CastContext)

	// EnterParenthesizedExpression is called when entering the parenthesizedExpression production.
	EnterParenthesizedExpression(c *ParenthesizedExpressionContext)

	// EnterUserVariable is called when entering the userVariable production.
	EnterUserVariable(c *UserVariableContext)

	// EnterElementAt is called when entering the elementAt production.
	EnterElementAt(c *ElementAtContext)

	// EnterLocalTimestamp is called when entering the localTimestamp production.
	EnterLocalTimestamp(c *LocalTimestampContext)

	// EnterCharFunction is called when entering the charFunction production.
	EnterCharFunction(c *CharFunctionContext)

	// EnterIntervalLiteral is called when entering the intervalLiteral production.
	EnterIntervalLiteral(c *IntervalLiteralContext)

	// EnterArrayRange is called when entering the arrayRange production.
	EnterArrayRange(c *ArrayRangeContext)

	// EnterSimpleCase is called when entering the simpleCase production.
	EnterSimpleCase(c *SimpleCaseContext)

	// EnterColumnReference is called when entering the columnReference production.
	EnterColumnReference(c *ColumnReferenceContext)

	// EnterStar is called when entering the star production.
	EnterStar(c *StarContext)

	// EnterConvertType is called when entering the convertType production.
	EnterConvertType(c *ConvertTypeContext)

	// EnterTimestampdiff is called when entering the timestampdiff production.
	EnterTimestampdiff(c *TimestampdiffContext)

	// EnterConvertCharSet is called when entering the convertCharSet production.
	EnterConvertCharSet(c *ConvertCharSetContext)

	// EnterSubqueryExpression is called when entering the subqueryExpression production.
	EnterSubqueryExpression(c *SubqueryExpressionContext)

	// EnterEncryptKey is called when entering the encryptKey production.
	EnterEncryptKey(c *EncryptKeyContext)

	// EnterDate_add is called when entering the date_add production.
	EnterDate_add(c *Date_addContext)

	// EnterCurrentTime is called when entering the currentTime production.
	EnterCurrentTime(c *CurrentTimeContext)

	// EnterLocalTime is called when entering the localTime production.
	EnterLocalTime(c *LocalTimeContext)

	// EnterSystemVariable is called when entering the systemVariable production.
	EnterSystemVariable(c *SystemVariableContext)

	// EnterCollate is called when entering the collate production.
	EnterCollate(c *CollateContext)

	// EnterCurrentUser is called when entering the currentUser production.
	EnterCurrentUser(c *CurrentUserContext)

	// EnterConstantDefault is called when entering the constantDefault production.
	EnterConstantDefault(c *ConstantDefaultContext)

	// EnterExtract is called when entering the extract production.
	EnterExtract(c *ExtractContext)

	// EnterCurrentTimestamp is called when entering the currentTimestamp production.
	EnterCurrentTimestamp(c *CurrentTimestampContext)

	// EnterFunctionCall is called when entering the functionCall production.
	EnterFunctionCall(c *FunctionCallContext)

	// EnterArraySlice is called when entering the arraySlice production.
	EnterArraySlice(c *ArraySliceContext)

	// EnterDateFloor is called when entering the dateFloor production.
	EnterDateFloor(c *DateFloorContext)

	// EnterSearchedCase is called when entering the searchedCase production.
	EnterSearchedCase(c *SearchedCaseContext)

	// EnterExcept is called when entering the except production.
	EnterExcept(c *ExceptContext)

	// EnterReplace is called when entering the replace production.
	EnterReplace(c *ReplaceContext)

	// EnterCastDataType is called when entering the castDataType production.
	EnterCastDataType(c *CastDataTypeContext)

	// EnterFunctionCallExpression is called when entering the functionCallExpression production.
	EnterFunctionCallExpression(c *FunctionCallExpressionContext)

	// EnterFunctionIdentifier is called when entering the functionIdentifier production.
	EnterFunctionIdentifier(c *FunctionIdentifierContext)

	// EnterFunctionNameIdentifier is called when entering the functionNameIdentifier production.
	EnterFunctionNameIdentifier(c *FunctionNameIdentifierContext)

	// EnterWindowSpec is called when entering the windowSpec production.
	EnterWindowSpec(c *WindowSpecContext)

	// EnterWindowFrame is called when entering the windowFrame production.
	EnterWindowFrame(c *WindowFrameContext)

	// EnterFrameUnits is called when entering the frameUnits production.
	EnterFrameUnits(c *FrameUnitsContext)

	// EnterFrameBoundary is called when entering the frameBoundary production.
	EnterFrameBoundary(c *FrameBoundaryContext)

	// EnterQualifiedName is called when entering the qualifiedName production.
	EnterQualifiedName(c *QualifiedNameContext)

	// EnterSpecifiedPartition is called when entering the specifiedPartition production.
	EnterSpecifiedPartition(c *SpecifiedPartitionContext)

	// EnterNullLiteral is called when entering the nullLiteral production.
	EnterNullLiteral(c *NullLiteralContext)

	// EnterTypeConstructor is called when entering the typeConstructor production.
	EnterTypeConstructor(c *TypeConstructorContext)

	// EnterNumericLiteral is called when entering the numericLiteral production.
	EnterNumericLiteral(c *NumericLiteralContext)

	// EnterBooleanLiteral is called when entering the booleanLiteral production.
	EnterBooleanLiteral(c *BooleanLiteralContext)

	// EnterStringLiteral is called when entering the stringLiteral production.
	EnterStringLiteral(c *StringLiteralContext)

	// EnterArrayLiteral is called when entering the arrayLiteral production.
	EnterArrayLiteral(c *ArrayLiteralContext)

	// EnterMapLiteral is called when entering the mapLiteral production.
	EnterMapLiteral(c *MapLiteralContext)

	// EnterStructLiteral is called when entering the structLiteral production.
	EnterStructLiteral(c *StructLiteralContext)

	// EnterPlaceholder is called when entering the placeholder production.
	EnterPlaceholder(c *PlaceholderContext)

	// EnterComparisonOperator is called when entering the comparisonOperator production.
	EnterComparisonOperator(c *ComparisonOperatorContext)

	// EnterBooleanValue is called when entering the booleanValue production.
	EnterBooleanValue(c *BooleanValueContext)

	// EnterWhenClause is called when entering the whenClause production.
	EnterWhenClause(c *WhenClauseContext)

	// EnterInterval is called when entering the interval production.
	EnterInterval(c *IntervalContext)

	// EnterUnitIdentifier is called when entering the unitIdentifier production.
	EnterUnitIdentifier(c *UnitIdentifierContext)

	// EnterDataTypeWithNullable is called when entering the dataTypeWithNullable production.
	EnterDataTypeWithNullable(c *DataTypeWithNullableContext)

	// EnterComplexDataType is called when entering the complexDataType production.
	EnterComplexDataType(c *ComplexDataTypeContext)

	// EnterAggStateDataType is called when entering the aggStateDataType production.
	EnterAggStateDataType(c *AggStateDataTypeContext)

	// EnterPrimitiveDataType is called when entering the primitiveDataType production.
	EnterPrimitiveDataType(c *PrimitiveDataTypeContext)

	// EnterPrimitiveColType is called when entering the primitiveColType production.
	EnterPrimitiveColType(c *PrimitiveColTypeContext)

	// EnterComplexColTypeList is called when entering the complexColTypeList production.
	EnterComplexColTypeList(c *ComplexColTypeListContext)

	// EnterComplexColType is called when entering the complexColType production.
	EnterComplexColType(c *ComplexColTypeContext)

	// EnterCommentSpec is called when entering the commentSpec production.
	EnterCommentSpec(c *CommentSpecContext)

	// EnterSample is called when entering the sample production.
	EnterSample(c *SampleContext)

	// EnterSampleByPercentile is called when entering the sampleByPercentile production.
	EnterSampleByPercentile(c *SampleByPercentileContext)

	// EnterSampleByRows is called when entering the sampleByRows production.
	EnterSampleByRows(c *SampleByRowsContext)

	// EnterTableSnapshot is called when entering the tableSnapshot production.
	EnterTableSnapshot(c *TableSnapshotContext)

	// EnterErrorCapturingIdentifier is called when entering the errorCapturingIdentifier production.
	EnterErrorCapturingIdentifier(c *ErrorCapturingIdentifierContext)

	// EnterErrorIdent is called when entering the errorIdent production.
	EnterErrorIdent(c *ErrorIdentContext)

	// EnterRealIdent is called when entering the realIdent production.
	EnterRealIdent(c *RealIdentContext)

	// EnterIdentifier is called when entering the identifier production.
	EnterIdentifier(c *IdentifierContext)

	// EnterUnquotedIdentifier is called when entering the unquotedIdentifier production.
	EnterUnquotedIdentifier(c *UnquotedIdentifierContext)

	// EnterQuotedIdentifierAlternative is called when entering the quotedIdentifierAlternative production.
	EnterQuotedIdentifierAlternative(c *QuotedIdentifierAlternativeContext)

	// EnterQuotedIdentifier is called when entering the quotedIdentifier production.
	EnterQuotedIdentifier(c *QuotedIdentifierContext)

	// EnterIntegerLiteral is called when entering the integerLiteral production.
	EnterIntegerLiteral(c *IntegerLiteralContext)

	// EnterDecimalLiteral is called when entering the decimalLiteral production.
	EnterDecimalLiteral(c *DecimalLiteralContext)

	// EnterNonReserved is called when entering the nonReserved production.
	EnterNonReserved(c *NonReservedContext)

	// ExitMultiStatements is called when exiting the multiStatements production.
	ExitMultiStatements(c *MultiStatementsContext)

	// ExitSingleStatement is called when exiting the singleStatement production.
	ExitSingleStatement(c *SingleStatementContext)

	// ExitStatementBaseAlias is called when exiting the statementBaseAlias production.
	ExitStatementBaseAlias(c *StatementBaseAliasContext)

	// ExitCallProcedure is called when exiting the callProcedure production.
	ExitCallProcedure(c *CallProcedureContext)

	// ExitCreateProcedure is called when exiting the createProcedure production.
	ExitCreateProcedure(c *CreateProcedureContext)

	// ExitDropProcedure is called when exiting the dropProcedure production.
	ExitDropProcedure(c *DropProcedureContext)

	// ExitShowProcedureStatus is called when exiting the showProcedureStatus production.
	ExitShowProcedureStatus(c *ShowProcedureStatusContext)

	// ExitShowCreateProcedure is called when exiting the showCreateProcedure production.
	ExitShowCreateProcedure(c *ShowCreateProcedureContext)

	// ExitShowConfig is called when exiting the showConfig production.
	ExitShowConfig(c *ShowConfigContext)

	// ExitStatementDefault is called when exiting the statementDefault production.
	ExitStatementDefault(c *StatementDefaultContext)

	// ExitSupportedDmlStatementAlias is called when exiting the supportedDmlStatementAlias production.
	ExitSupportedDmlStatementAlias(c *SupportedDmlStatementAliasContext)

	// ExitSupportedCreateStatementAlias is called when exiting the supportedCreateStatementAlias production.
	ExitSupportedCreateStatementAlias(c *SupportedCreateStatementAliasContext)

	// ExitSupportedAlterStatementAlias is called when exiting the supportedAlterStatementAlias production.
	ExitSupportedAlterStatementAlias(c *SupportedAlterStatementAliasContext)

	// ExitMaterailizedViewStatementAlias is called when exiting the materailizedViewStatementAlias production.
	ExitMaterailizedViewStatementAlias(c *MaterailizedViewStatementAliasContext)

	// ExitConstraintStatementAlias is called when exiting the constraintStatementAlias production.
	ExitConstraintStatementAlias(c *ConstraintStatementAliasContext)

	// ExitSupportedDropStatementAlias is called when exiting the supportedDropStatementAlias production.
	ExitSupportedDropStatementAlias(c *SupportedDropStatementAliasContext)

	// ExitUnsupported is called when exiting the unsupported production.
	ExitUnsupported(c *UnsupportedContext)

	// ExitUnsupportedStatement is called when exiting the unsupportedStatement production.
	ExitUnsupportedStatement(c *UnsupportedStatementContext)

	// ExitCreateMTMV is called when exiting the createMTMV production.
	ExitCreateMTMV(c *CreateMTMVContext)

	// ExitRefreshMTMV is called when exiting the refreshMTMV production.
	ExitRefreshMTMV(c *RefreshMTMVContext)

	// ExitAlterMTMV is called when exiting the alterMTMV production.
	ExitAlterMTMV(c *AlterMTMVContext)

	// ExitDropMTMV is called when exiting the dropMTMV production.
	ExitDropMTMV(c *DropMTMVContext)

	// ExitPauseMTMV is called when exiting the pauseMTMV production.
	ExitPauseMTMV(c *PauseMTMVContext)

	// ExitResumeMTMV is called when exiting the resumeMTMV production.
	ExitResumeMTMV(c *ResumeMTMVContext)

	// ExitCancelMTMVTask is called when exiting the cancelMTMVTask production.
	ExitCancelMTMVTask(c *CancelMTMVTaskContext)

	// ExitShowCreateMTMV is called when exiting the showCreateMTMV production.
	ExitShowCreateMTMV(c *ShowCreateMTMVContext)

	// ExitAddConstraint is called when exiting the addConstraint production.
	ExitAddConstraint(c *AddConstraintContext)

	// ExitDropConstraint is called when exiting the dropConstraint production.
	ExitDropConstraint(c *DropConstraintContext)

	// ExitShowConstraint is called when exiting the showConstraint production.
	ExitShowConstraint(c *ShowConstraintContext)

	// ExitInsertTable is called when exiting the insertTable production.
	ExitInsertTable(c *InsertTableContext)

	// ExitUpdate is called when exiting the update production.
	ExitUpdate(c *UpdateContext)

	// ExitDelete is called when exiting the delete production.
	ExitDelete(c *DeleteContext)

	// ExitLoad is called when exiting the load production.
	ExitLoad(c *LoadContext)

	// ExitMysqlLoad is called when exiting the mysqlLoad production.
	ExitMysqlLoad(c *MysqlLoadContext)

	// ExitExport is called when exiting the export production.
	ExitExport(c *ExportContext)

	// ExitCreateTable is called when exiting the createTable production.
	ExitCreateTable(c *CreateTableContext)

	// ExitCreateView is called when exiting the createView production.
	ExitCreateView(c *CreateViewContext)

	// ExitCreateTableLike is called when exiting the createTableLike production.
	ExitCreateTableLike(c *CreateTableLikeContext)

	// ExitCreateRowPolicy is called when exiting the createRowPolicy production.
	ExitCreateRowPolicy(c *CreateRowPolicyContext)

	// ExitAlterView is called when exiting the alterView production.
	ExitAlterView(c *AlterViewContext)

	// ExitAlterStorageVault is called when exiting the alterStorageVault production.
	ExitAlterStorageVault(c *AlterStorageVaultContext)

	// ExitDropCatalogRecycleBin is called when exiting the dropCatalogRecycleBin production.
	ExitDropCatalogRecycleBin(c *DropCatalogRecycleBinContext)

	// ExitAdminShowReplicaStatus is called when exiting the adminShowReplicaStatus production.
	ExitAdminShowReplicaStatus(c *AdminShowReplicaStatusContext)

	// ExitAdminShowReplicaDistribution is called when exiting the adminShowReplicaDistribution production.
	ExitAdminShowReplicaDistribution(c *AdminShowReplicaDistributionContext)

	// ExitAdminSetReplicaStatus is called when exiting the adminSetReplicaStatus production.
	ExitAdminSetReplicaStatus(c *AdminSetReplicaStatusContext)

	// ExitAdminSetReplicaVersion is called when exiting the adminSetReplicaVersion production.
	ExitAdminSetReplicaVersion(c *AdminSetReplicaVersionContext)

	// ExitAdminRepairTable is called when exiting the adminRepairTable production.
	ExitAdminRepairTable(c *AdminRepairTableContext)

	// ExitAdminCancelRepairTable is called when exiting the adminCancelRepairTable production.
	ExitAdminCancelRepairTable(c *AdminCancelRepairTableContext)

	// ExitAdminCompactTable is called when exiting the adminCompactTable production.
	ExitAdminCompactTable(c *AdminCompactTableContext)

	// ExitAdminSetFrontendConfig is called when exiting the adminSetFrontendConfig production.
	ExitAdminSetFrontendConfig(c *AdminSetFrontendConfigContext)

	// ExitAdminCheckTablets is called when exiting the adminCheckTablets production.
	ExitAdminCheckTablets(c *AdminCheckTabletsContext)

	// ExitAdminRebalanceDisk is called when exiting the adminRebalanceDisk production.
	ExitAdminRebalanceDisk(c *AdminRebalanceDiskContext)

	// ExitAdminCancelRebalanceDisk is called when exiting the adminCancelRebalanceDisk production.
	ExitAdminCancelRebalanceDisk(c *AdminCancelRebalanceDiskContext)

	// ExitAdminCleanTrash is called when exiting the adminCleanTrash production.
	ExitAdminCleanTrash(c *AdminCleanTrashContext)

	// ExitAdminSetPartitionVersion is called when exiting the adminSetPartitionVersion production.
	ExitAdminSetPartitionVersion(c *AdminSetPartitionVersionContext)

	// ExitAdminDiagnoseTablet is called when exiting the adminDiagnoseTablet production.
	ExitAdminDiagnoseTablet(c *AdminDiagnoseTabletContext)

	// ExitAdminShowTabletStorageFormat is called when exiting the adminShowTabletStorageFormat production.
	ExitAdminShowTabletStorageFormat(c *AdminShowTabletStorageFormatContext)

	// ExitAdminCopyTablet is called when exiting the adminCopyTablet production.
	ExitAdminCopyTablet(c *AdminCopyTabletContext)

	// ExitAdminSetTableStatus is called when exiting the adminSetTableStatus production.
	ExitAdminSetTableStatus(c *AdminSetTableStatusContext)

	// ExitBaseTableRef is called when exiting the baseTableRef production.
	ExitBaseTableRef(c *BaseTableRefContext)

	// ExitWildWhere is called when exiting the wildWhere production.
	ExitWildWhere(c *WildWhereContext)

	// ExitTransactionBegin is called when exiting the transactionBegin production.
	ExitTransactionBegin(c *TransactionBeginContext)

	// ExitTranscationCommit is called when exiting the transcationCommit production.
	ExitTranscationCommit(c *TranscationCommitContext)

	// ExitTransactionRollback is called when exiting the transactionRollback production.
	ExitTransactionRollback(c *TransactionRollbackContext)

	// ExitGrantTablePrivilege is called when exiting the grantTablePrivilege production.
	ExitGrantTablePrivilege(c *GrantTablePrivilegeContext)

	// ExitGrantResourcePrivilege is called when exiting the grantResourcePrivilege production.
	ExitGrantResourcePrivilege(c *GrantResourcePrivilegeContext)

	// ExitGrantRole is called when exiting the grantRole production.
	ExitGrantRole(c *GrantRoleContext)

	// ExitPrivilege is called when exiting the privilege production.
	ExitPrivilege(c *PrivilegeContext)

	// ExitPrivilegeList is called when exiting the privilegeList production.
	ExitPrivilegeList(c *PrivilegeListContext)

	// ExitAlterTable is called when exiting the alterTable production.
	ExitAlterTable(c *AlterTableContext)

	// ExitAlterTableAddRollup is called when exiting the alterTableAddRollup production.
	ExitAlterTableAddRollup(c *AlterTableAddRollupContext)

	// ExitAlterTableDropRollup is called when exiting the alterTableDropRollup production.
	ExitAlterTableDropRollup(c *AlterTableDropRollupContext)

	// ExitAlterSystem is called when exiting the alterSystem production.
	ExitAlterSystem(c *AlterSystemContext)

	// ExitAlterDatabaseSetQuota is called when exiting the alterDatabaseSetQuota production.
	ExitAlterDatabaseSetQuota(c *AlterDatabaseSetQuotaContext)

	// ExitAlterDatabaseRename is called when exiting the alterDatabaseRename production.
	ExitAlterDatabaseRename(c *AlterDatabaseRenameContext)

	// ExitAlterDatabaseProperties is called when exiting the alterDatabaseProperties production.
	ExitAlterDatabaseProperties(c *AlterDatabasePropertiesContext)

	// ExitAlterCatalogRename is called when exiting the alterCatalogRename production.
	ExitAlterCatalogRename(c *AlterCatalogRenameContext)

	// ExitAlterCatalogProperties is called when exiting the alterCatalogProperties production.
	ExitAlterCatalogProperties(c *AlterCatalogPropertiesContext)

	// ExitAlterCatalogComment is called when exiting the alterCatalogComment production.
	ExitAlterCatalogComment(c *AlterCatalogCommentContext)

	// ExitAlterResource is called when exiting the alterResource production.
	ExitAlterResource(c *AlterResourceContext)

	// ExitAlterColocateGroup is called when exiting the alterColocateGroup production.
	ExitAlterColocateGroup(c *AlterColocateGroupContext)

	// ExitAlterWorkloadGroup is called when exiting the alterWorkloadGroup production.
	ExitAlterWorkloadGroup(c *AlterWorkloadGroupContext)

	// ExitAlterWorkloadPolicy is called when exiting the alterWorkloadPolicy production.
	ExitAlterWorkloadPolicy(c *AlterWorkloadPolicyContext)

	// ExitAlterRoutineLoad is called when exiting the alterRoutineLoad production.
	ExitAlterRoutineLoad(c *AlterRoutineLoadContext)

	// ExitAlterSqlBlockRule is called when exiting the alterSqlBlockRule production.
	ExitAlterSqlBlockRule(c *AlterSqlBlockRuleContext)

	// ExitAlterTableProperties is called when exiting the alterTableProperties production.
	ExitAlterTableProperties(c *AlterTablePropertiesContext)

	// ExitAlterStoragePlicy is called when exiting the alterStoragePlicy production.
	ExitAlterStoragePlicy(c *AlterStoragePlicyContext)

	// ExitAlterUser is called when exiting the alterUser production.
	ExitAlterUser(c *AlterUserContext)

	// ExitAlterRepository is called when exiting the alterRepository production.
	ExitAlterRepository(c *AlterRepositoryContext)

	// ExitAddBackendClause is called when exiting the addBackendClause production.
	ExitAddBackendClause(c *AddBackendClauseContext)

	// ExitDropBackendClause is called when exiting the dropBackendClause production.
	ExitDropBackendClause(c *DropBackendClauseContext)

	// ExitDecommissionBackendClause is called when exiting the decommissionBackendClause production.
	ExitDecommissionBackendClause(c *DecommissionBackendClauseContext)

	// ExitAddObserverClause is called when exiting the addObserverClause production.
	ExitAddObserverClause(c *AddObserverClauseContext)

	// ExitDropObserverClause is called when exiting the dropObserverClause production.
	ExitDropObserverClause(c *DropObserverClauseContext)

	// ExitAddFollowerClause is called when exiting the addFollowerClause production.
	ExitAddFollowerClause(c *AddFollowerClauseContext)

	// ExitDropFollowerClause is called when exiting the dropFollowerClause production.
	ExitDropFollowerClause(c *DropFollowerClauseContext)

	// ExitAddBrokerClause is called when exiting the addBrokerClause production.
	ExitAddBrokerClause(c *AddBrokerClauseContext)

	// ExitDropBrokerClause is called when exiting the dropBrokerClause production.
	ExitDropBrokerClause(c *DropBrokerClauseContext)

	// ExitDropAllBrokerClause is called when exiting the dropAllBrokerClause production.
	ExitDropAllBrokerClause(c *DropAllBrokerClauseContext)

	// ExitAlterLoadErrorUrlClause is called when exiting the alterLoadErrorUrlClause production.
	ExitAlterLoadErrorUrlClause(c *AlterLoadErrorUrlClauseContext)

	// ExitModifyBackendClause is called when exiting the modifyBackendClause production.
	ExitModifyBackendClause(c *ModifyBackendClauseContext)

	// ExitModifyFrontendOrBackendHostNameClause is called when exiting the modifyFrontendOrBackendHostNameClause production.
	ExitModifyFrontendOrBackendHostNameClause(c *ModifyFrontendOrBackendHostNameClauseContext)

	// ExitDropRollupClause is called when exiting the dropRollupClause production.
	ExitDropRollupClause(c *DropRollupClauseContext)

	// ExitAddRollupClause is called when exiting the addRollupClause production.
	ExitAddRollupClause(c *AddRollupClauseContext)

	// ExitAddColumnClause is called when exiting the addColumnClause production.
	ExitAddColumnClause(c *AddColumnClauseContext)

	// ExitAddColumnsClause is called when exiting the addColumnsClause production.
	ExitAddColumnsClause(c *AddColumnsClauseContext)

	// ExitDropColumnClause is called when exiting the dropColumnClause production.
	ExitDropColumnClause(c *DropColumnClauseContext)

	// ExitModifyColumnClause is called when exiting the modifyColumnClause production.
	ExitModifyColumnClause(c *ModifyColumnClauseContext)

	// ExitReorderColumnsClause is called when exiting the reorderColumnsClause production.
	ExitReorderColumnsClause(c *ReorderColumnsClauseContext)

	// ExitAddPartitionClause is called when exiting the addPartitionClause production.
	ExitAddPartitionClause(c *AddPartitionClauseContext)

	// ExitDropPartitionClause is called when exiting the dropPartitionClause production.
	ExitDropPartitionClause(c *DropPartitionClauseContext)

	// ExitModifyPartitionClause is called when exiting the modifyPartitionClause production.
	ExitModifyPartitionClause(c *ModifyPartitionClauseContext)

	// ExitReplacePartitionClause is called when exiting the replacePartitionClause production.
	ExitReplacePartitionClause(c *ReplacePartitionClauseContext)

	// ExitReplaceTableClause is called when exiting the replaceTableClause production.
	ExitReplaceTableClause(c *ReplaceTableClauseContext)

	// ExitRenameClause is called when exiting the renameClause production.
	ExitRenameClause(c *RenameClauseContext)

	// ExitRenameRollupClause is called when exiting the renameRollupClause production.
	ExitRenameRollupClause(c *RenameRollupClauseContext)

	// ExitRenamePartitionClause is called when exiting the renamePartitionClause production.
	ExitRenamePartitionClause(c *RenamePartitionClauseContext)

	// ExitRenameColumnClause is called when exiting the renameColumnClause production.
	ExitRenameColumnClause(c *RenameColumnClauseContext)

	// ExitAddIndexClause is called when exiting the addIndexClause production.
	ExitAddIndexClause(c *AddIndexClauseContext)

	// ExitDropIndexClause is called when exiting the dropIndexClause production.
	ExitDropIndexClause(c *DropIndexClauseContext)

	// ExitEnableFeatureClause is called when exiting the enableFeatureClause production.
	ExitEnableFeatureClause(c *EnableFeatureClauseContext)

	// ExitModifyDistributionClause is called when exiting the modifyDistributionClause production.
	ExitModifyDistributionClause(c *ModifyDistributionClauseContext)

	// ExitModifyTableCommentClause is called when exiting the modifyTableCommentClause production.
	ExitModifyTableCommentClause(c *ModifyTableCommentClauseContext)

	// ExitModifyColumnCommentClause is called when exiting the modifyColumnCommentClause production.
	ExitModifyColumnCommentClause(c *ModifyColumnCommentClauseContext)

	// ExitModifyEngineClause is called when exiting the modifyEngineClause production.
	ExitModifyEngineClause(c *ModifyEngineClauseContext)

	// ExitAlterMultiPartitionClause is called when exiting the alterMultiPartitionClause production.
	ExitAlterMultiPartitionClause(c *AlterMultiPartitionClauseContext)

	// ExitColumnPosition is called when exiting the columnPosition production.
	ExitColumnPosition(c *ColumnPositionContext)

	// ExitToRollup is called when exiting the toRollup production.
	ExitToRollup(c *ToRollupContext)

	// ExitFromRollup is called when exiting the fromRollup production.
	ExitFromRollup(c *FromRollupContext)

	// ExitDropDatabase is called when exiting the dropDatabase production.
	ExitDropDatabase(c *DropDatabaseContext)

	// ExitDropCatalog is called when exiting the dropCatalog production.
	ExitDropCatalog(c *DropCatalogContext)

	// ExitDropFunction is called when exiting the dropFunction production.
	ExitDropFunction(c *DropFunctionContext)

	// ExitDropTable is called when exiting the dropTable production.
	ExitDropTable(c *DropTableContext)

	// ExitDropUser is called when exiting the dropUser production.
	ExitDropUser(c *DropUserContext)

	// ExitDropView is called when exiting the dropView production.
	ExitDropView(c *DropViewContext)

	// ExitDropRepository is called when exiting the dropRepository production.
	ExitDropRepository(c *DropRepositoryContext)

	// ExitDropRole is called when exiting the dropRole production.
	ExitDropRole(c *DropRoleContext)

	// ExitDropFile is called when exiting the dropFile production.
	ExitDropFile(c *DropFileContext)

	// ExitDropIndex is called when exiting the dropIndex production.
	ExitDropIndex(c *DropIndexContext)

	// ExitDropResource is called when exiting the dropResource production.
	ExitDropResource(c *DropResourceContext)

	// ExitDropWorkloadGroup is called when exiting the dropWorkloadGroup production.
	ExitDropWorkloadGroup(c *DropWorkloadGroupContext)

	// ExitDropWorkloadPolicy is called when exiting the dropWorkloadPolicy production.
	ExitDropWorkloadPolicy(c *DropWorkloadPolicyContext)

	// ExitDropEncryptkey is called when exiting the dropEncryptkey production.
	ExitDropEncryptkey(c *DropEncryptkeyContext)

	// ExitDropSqlBlockRule is called when exiting the dropSqlBlockRule production.
	ExitDropSqlBlockRule(c *DropSqlBlockRuleContext)

	// ExitDropRowPolicy is called when exiting the dropRowPolicy production.
	ExitDropRowPolicy(c *DropRowPolicyContext)

	// ExitDropStoragePolicy is called when exiting the dropStoragePolicy production.
	ExitDropStoragePolicy(c *DropStoragePolicyContext)

	// ExitDropStage is called when exiting the dropStage production.
	ExitDropStage(c *DropStageContext)

	// ExitAlterTableStats is called when exiting the alterTableStats production.
	ExitAlterTableStats(c *AlterTableStatsContext)

	// ExitAlterColumnStats is called when exiting the alterColumnStats production.
	ExitAlterColumnStats(c *AlterColumnStatsContext)

	// ExitDropStats is called when exiting the dropStats production.
	ExitDropStats(c *DropStatsContext)

	// ExitDropCachedStats is called when exiting the dropCachedStats production.
	ExitDropCachedStats(c *DropCachedStatsContext)

	// ExitDropExpiredStats is called when exiting the dropExpiredStats production.
	ExitDropExpiredStats(c *DropExpiredStatsContext)

	// ExitDropAanalyzeJob is called when exiting the dropAanalyzeJob production.
	ExitDropAanalyzeJob(c *DropAanalyzeJobContext)

	// ExitCreateDatabase is called when exiting the createDatabase production.
	ExitCreateDatabase(c *CreateDatabaseContext)

	// ExitCreateCatalog is called when exiting the createCatalog production.
	ExitCreateCatalog(c *CreateCatalogContext)

	// ExitCreateUserDefineFunction is called when exiting the createUserDefineFunction production.
	ExitCreateUserDefineFunction(c *CreateUserDefineFunctionContext)

	// ExitCreateAliasFunction is called when exiting the createAliasFunction production.
	ExitCreateAliasFunction(c *CreateAliasFunctionContext)

	// ExitCreateUser is called when exiting the createUser production.
	ExitCreateUser(c *CreateUserContext)

	// ExitCreateRepository is called when exiting the createRepository production.
	ExitCreateRepository(c *CreateRepositoryContext)

	// ExitCreateRole is called when exiting the createRole production.
	ExitCreateRole(c *CreateRoleContext)

	// ExitCreateFile is called when exiting the createFile production.
	ExitCreateFile(c *CreateFileContext)

	// ExitCreateIndex is called when exiting the createIndex production.
	ExitCreateIndex(c *CreateIndexContext)

	// ExitCreateResource is called when exiting the createResource production.
	ExitCreateResource(c *CreateResourceContext)

	// ExitCreateStorageVault is called when exiting the createStorageVault production.
	ExitCreateStorageVault(c *CreateStorageVaultContext)

	// ExitCreateWorkloadGroup is called when exiting the createWorkloadGroup production.
	ExitCreateWorkloadGroup(c *CreateWorkloadGroupContext)

	// ExitCreateWorkloadPolicy is called when exiting the createWorkloadPolicy production.
	ExitCreateWorkloadPolicy(c *CreateWorkloadPolicyContext)

	// ExitCreateEncryptkey is called when exiting the createEncryptkey production.
	ExitCreateEncryptkey(c *CreateEncryptkeyContext)

	// ExitCreateDataSyncJob is called when exiting the createDataSyncJob production.
	ExitCreateDataSyncJob(c *CreateDataSyncJobContext)

	// ExitCreateSqlBlockRule is called when exiting the createSqlBlockRule production.
	ExitCreateSqlBlockRule(c *CreateSqlBlockRuleContext)

	// ExitCreateStoragePolicy is called when exiting the createStoragePolicy production.
	ExitCreateStoragePolicy(c *CreateStoragePolicyContext)

	// ExitBuildIndex is called when exiting the buildIndex production.
	ExitBuildIndex(c *BuildIndexContext)

	// ExitCreateStage is called when exiting the createStage production.
	ExitCreateStage(c *CreateStageContext)

	// ExitChannelDescriptions is called when exiting the channelDescriptions production.
	ExitChannelDescriptions(c *ChannelDescriptionsContext)

	// ExitChannelDescription is called when exiting the channelDescription production.
	ExitChannelDescription(c *ChannelDescriptionContext)

	// ExitWorkloadPolicyActions is called when exiting the workloadPolicyActions production.
	ExitWorkloadPolicyActions(c *WorkloadPolicyActionsContext)

	// ExitWorkloadPolicyAction is called when exiting the workloadPolicyAction production.
	ExitWorkloadPolicyAction(c *WorkloadPolicyActionContext)

	// ExitWorkloadPolicyConditions is called when exiting the workloadPolicyConditions production.
	ExitWorkloadPolicyConditions(c *WorkloadPolicyConditionsContext)

	// ExitWorkloadPolicyCondition is called when exiting the workloadPolicyCondition production.
	ExitWorkloadPolicyCondition(c *WorkloadPolicyConditionContext)

	// ExitStorageBackend is called when exiting the storageBackend production.
	ExitStorageBackend(c *StorageBackendContext)

	// ExitPasswordOption is called when exiting the passwordOption production.
	ExitPasswordOption(c *PasswordOptionContext)

	// ExitFunctionArguments is called when exiting the functionArguments production.
	ExitFunctionArguments(c *FunctionArgumentsContext)

	// ExitFunctionArgument is called when exiting the functionArgument production.
	ExitFunctionArgument(c *FunctionArgumentContext)

	// ExitSetOptions is called when exiting the setOptions production.
	ExitSetOptions(c *SetOptionsContext)

	// ExitSetDefaultStorageVault is called when exiting the setDefaultStorageVault production.
	ExitSetDefaultStorageVault(c *SetDefaultStorageVaultContext)

	// ExitSetUserProperties is called when exiting the setUserProperties production.
	ExitSetUserProperties(c *SetUserPropertiesContext)

	// ExitSetTransaction is called when exiting the setTransaction production.
	ExitSetTransaction(c *SetTransactionContext)

	// ExitOptionWithType is called when exiting the optionWithType production.
	ExitOptionWithType(c *OptionWithTypeContext)

	// ExitSetNames is called when exiting the setNames production.
	ExitSetNames(c *SetNamesContext)

	// ExitSetCharset is called when exiting the setCharset production.
	ExitSetCharset(c *SetCharsetContext)

	// ExitSetCollate is called when exiting the setCollate production.
	ExitSetCollate(c *SetCollateContext)

	// ExitSetPassword is called when exiting the setPassword production.
	ExitSetPassword(c *SetPasswordContext)

	// ExitSetLdapAdminPassword is called when exiting the setLdapAdminPassword production.
	ExitSetLdapAdminPassword(c *SetLdapAdminPasswordContext)

	// ExitSetVariableWithoutType is called when exiting the setVariableWithoutType production.
	ExitSetVariableWithoutType(c *SetVariableWithoutTypeContext)

	// ExitSetSystemVariable is called when exiting the setSystemVariable production.
	ExitSetSystemVariable(c *SetSystemVariableContext)

	// ExitSetUserVariable is called when exiting the setUserVariable production.
	ExitSetUserVariable(c *SetUserVariableContext)

	// ExitTransactionAccessMode is called when exiting the transactionAccessMode production.
	ExitTransactionAccessMode(c *TransactionAccessModeContext)

	// ExitIsolationLevel is called when exiting the isolationLevel production.
	ExitIsolationLevel(c *IsolationLevelContext)

	// ExitUnsupoortedUnsetStatement is called when exiting the unsupoortedUnsetStatement production.
	ExitUnsupoortedUnsetStatement(c *UnsupoortedUnsetStatementContext)

	// ExitUseDatabase is called when exiting the useDatabase production.
	ExitUseDatabase(c *UseDatabaseContext)

	// ExitUseCloudCluster is called when exiting the useCloudCluster production.
	ExitUseCloudCluster(c *UseCloudClusterContext)

	// ExitSwitchCatalog is called when exiting the switchCatalog production.
	ExitSwitchCatalog(c *SwitchCatalogContext)

	// ExitTruncateTable is called when exiting the truncateTable production.
	ExitTruncateTable(c *TruncateTableContext)

	// ExitKillConnection is called when exiting the killConnection production.
	ExitKillConnection(c *KillConnectionContext)

	// ExitKillQuery is called when exiting the killQuery production.
	ExitKillQuery(c *KillQueryContext)

	// ExitDescribeTableValuedFunction is called when exiting the describeTableValuedFunction production.
	ExitDescribeTableValuedFunction(c *DescribeTableValuedFunctionContext)

	// ExitDescribeTableAll is called when exiting the describeTableAll production.
	ExitDescribeTableAll(c *DescribeTableAllContext)

	// ExitDescribeTable is called when exiting the describeTable production.
	ExitDescribeTable(c *DescribeTableContext)

	// ExitConstraint is called when exiting the constraint production.
	ExitConstraint(c *ConstraintContext)

	// ExitPartitionSpec is called when exiting the partitionSpec production.
	ExitPartitionSpec(c *PartitionSpecContext)

	// ExitPartitionTable is called when exiting the partitionTable production.
	ExitPartitionTable(c *PartitionTableContext)

	// ExitIdentityOrFunctionList is called when exiting the identityOrFunctionList production.
	ExitIdentityOrFunctionList(c *IdentityOrFunctionListContext)

	// ExitIdentityOrFunction is called when exiting the identityOrFunction production.
	ExitIdentityOrFunction(c *IdentityOrFunctionContext)

	// ExitDataDesc is called when exiting the dataDesc production.
	ExitDataDesc(c *DataDescContext)

	// ExitBuildMode is called when exiting the buildMode production.
	ExitBuildMode(c *BuildModeContext)

	// ExitRefreshTrigger is called when exiting the refreshTrigger production.
	ExitRefreshTrigger(c *RefreshTriggerContext)

	// ExitRefreshSchedule is called when exiting the refreshSchedule production.
	ExitRefreshSchedule(c *RefreshScheduleContext)

	// ExitRefreshMethod is called when exiting the refreshMethod production.
	ExitRefreshMethod(c *RefreshMethodContext)

	// ExitMvPartition is called when exiting the mvPartition production.
	ExitMvPartition(c *MvPartitionContext)

	// ExitIdentifierOrText is called when exiting the identifierOrText production.
	ExitIdentifierOrText(c *IdentifierOrTextContext)

	// ExitIdentifierOrTextOrAsterisk is called when exiting the identifierOrTextOrAsterisk production.
	ExitIdentifierOrTextOrAsterisk(c *IdentifierOrTextOrAsteriskContext)

	// ExitMultipartIdentifierOrAsterisk is called when exiting the multipartIdentifierOrAsterisk production.
	ExitMultipartIdentifierOrAsterisk(c *MultipartIdentifierOrAsteriskContext)

	// ExitIdentifierOrAsterisk is called when exiting the identifierOrAsterisk production.
	ExitIdentifierOrAsterisk(c *IdentifierOrAsteriskContext)

	// ExitUserIdentify is called when exiting the userIdentify production.
	ExitUserIdentify(c *UserIdentifyContext)

	// ExitGrantUserIdentify is called when exiting the grantUserIdentify production.
	ExitGrantUserIdentify(c *GrantUserIdentifyContext)

	// ExitExplain is called when exiting the explain production.
	ExitExplain(c *ExplainContext)

	// ExitExplainCommand is called when exiting the explainCommand production.
	ExitExplainCommand(c *ExplainCommandContext)

	// ExitPlanType is called when exiting the planType production.
	ExitPlanType(c *PlanTypeContext)

	// ExitMergeType is called when exiting the mergeType production.
	ExitMergeType(c *MergeTypeContext)

	// ExitPreFilterClause is called when exiting the preFilterClause production.
	ExitPreFilterClause(c *PreFilterClauseContext)

	// ExitDeleteOnClause is called when exiting the deleteOnClause production.
	ExitDeleteOnClause(c *DeleteOnClauseContext)

	// ExitSequenceColClause is called when exiting the sequenceColClause production.
	ExitSequenceColClause(c *SequenceColClauseContext)

	// ExitColFromPath is called when exiting the colFromPath production.
	ExitColFromPath(c *ColFromPathContext)

	// ExitColMappingList is called when exiting the colMappingList production.
	ExitColMappingList(c *ColMappingListContext)

	// ExitMappingExpr is called when exiting the mappingExpr production.
	ExitMappingExpr(c *MappingExprContext)

	// ExitWithRemoteStorageSystem is called when exiting the withRemoteStorageSystem production.
	ExitWithRemoteStorageSystem(c *WithRemoteStorageSystemContext)

	// ExitResourceDesc is called when exiting the resourceDesc production.
	ExitResourceDesc(c *ResourceDescContext)

	// ExitMysqlDataDesc is called when exiting the mysqlDataDesc production.
	ExitMysqlDataDesc(c *MysqlDataDescContext)

	// ExitSkipLines is called when exiting the skipLines production.
	ExitSkipLines(c *SkipLinesContext)

	// ExitOutFileClause is called when exiting the outFileClause production.
	ExitOutFileClause(c *OutFileClauseContext)

	// ExitQuery is called when exiting the query production.
	ExitQuery(c *QueryContext)

	// ExitQueryTermDefault is called when exiting the queryTermDefault production.
	ExitQueryTermDefault(c *QueryTermDefaultContext)

	// ExitSetOperation is called when exiting the setOperation production.
	ExitSetOperation(c *SetOperationContext)

	// ExitSetQuantifier is called when exiting the setQuantifier production.
	ExitSetQuantifier(c *SetQuantifierContext)

	// ExitQueryPrimaryDefault is called when exiting the queryPrimaryDefault production.
	ExitQueryPrimaryDefault(c *QueryPrimaryDefaultContext)

	// ExitSubquery is called when exiting the subquery production.
	ExitSubquery(c *SubqueryContext)

	// ExitValuesTable is called when exiting the valuesTable production.
	ExitValuesTable(c *ValuesTableContext)

	// ExitRegularQuerySpecification is called when exiting the regularQuerySpecification production.
	ExitRegularQuerySpecification(c *RegularQuerySpecificationContext)

	// ExitCte is called when exiting the cte production.
	ExitCte(c *CteContext)

	// ExitAliasQuery is called when exiting the aliasQuery production.
	ExitAliasQuery(c *AliasQueryContext)

	// ExitColumnAliases is called when exiting the columnAliases production.
	ExitColumnAliases(c *ColumnAliasesContext)

	// ExitSelectClause is called when exiting the selectClause production.
	ExitSelectClause(c *SelectClauseContext)

	// ExitSelectColumnClause is called when exiting the selectColumnClause production.
	ExitSelectColumnClause(c *SelectColumnClauseContext)

	// ExitWhereClause is called when exiting the whereClause production.
	ExitWhereClause(c *WhereClauseContext)

	// ExitFromClause is called when exiting the fromClause production.
	ExitFromClause(c *FromClauseContext)

	// ExitIntoClause is called when exiting the intoClause production.
	ExitIntoClause(c *IntoClauseContext)

	// ExitBulkCollectClause is called when exiting the bulkCollectClause production.
	ExitBulkCollectClause(c *BulkCollectClauseContext)

	// ExitTableRow is called when exiting the tableRow production.
	ExitTableRow(c *TableRowContext)

	// ExitRelations is called when exiting the relations production.
	ExitRelations(c *RelationsContext)

	// ExitRelation is called when exiting the relation production.
	ExitRelation(c *RelationContext)

	// ExitJoinRelation is called when exiting the joinRelation production.
	ExitJoinRelation(c *JoinRelationContext)

	// ExitBracketDistributeType is called when exiting the bracketDistributeType production.
	ExitBracketDistributeType(c *BracketDistributeTypeContext)

	// ExitCommentDistributeType is called when exiting the commentDistributeType production.
	ExitCommentDistributeType(c *CommentDistributeTypeContext)

	// ExitBracketRelationHint is called when exiting the bracketRelationHint production.
	ExitBracketRelationHint(c *BracketRelationHintContext)

	// ExitCommentRelationHint is called when exiting the commentRelationHint production.
	ExitCommentRelationHint(c *CommentRelationHintContext)

	// ExitAggClause is called when exiting the aggClause production.
	ExitAggClause(c *AggClauseContext)

	// ExitGroupingElement is called when exiting the groupingElement production.
	ExitGroupingElement(c *GroupingElementContext)

	// ExitGroupingSet is called when exiting the groupingSet production.
	ExitGroupingSet(c *GroupingSetContext)

	// ExitHavingClause is called when exiting the havingClause production.
	ExitHavingClause(c *HavingClauseContext)

	// ExitSelectHint is called when exiting the selectHint production.
	ExitSelectHint(c *SelectHintContext)

	// ExitHintStatement is called when exiting the hintStatement production.
	ExitHintStatement(c *HintStatementContext)

	// ExitHintAssignment is called when exiting the hintAssignment production.
	ExitHintAssignment(c *HintAssignmentContext)

	// ExitUpdateAssignment is called when exiting the updateAssignment production.
	ExitUpdateAssignment(c *UpdateAssignmentContext)

	// ExitUpdateAssignmentSeq is called when exiting the updateAssignmentSeq production.
	ExitUpdateAssignmentSeq(c *UpdateAssignmentSeqContext)

	// ExitLateralView is called when exiting the lateralView production.
	ExitLateralView(c *LateralViewContext)

	// ExitQueryOrganization is called when exiting the queryOrganization production.
	ExitQueryOrganization(c *QueryOrganizationContext)

	// ExitSortClause is called when exiting the sortClause production.
	ExitSortClause(c *SortClauseContext)

	// ExitSortItem is called when exiting the sortItem production.
	ExitSortItem(c *SortItemContext)

	// ExitLimitClause is called when exiting the limitClause production.
	ExitLimitClause(c *LimitClauseContext)

	// ExitPartitionClause is called when exiting the partitionClause production.
	ExitPartitionClause(c *PartitionClauseContext)

	// ExitJoinType is called when exiting the joinType production.
	ExitJoinType(c *JoinTypeContext)

	// ExitJoinCriteria is called when exiting the joinCriteria production.
	ExitJoinCriteria(c *JoinCriteriaContext)

	// ExitIdentifierList is called when exiting the identifierList production.
	ExitIdentifierList(c *IdentifierListContext)

	// ExitIdentifierSeq is called when exiting the identifierSeq production.
	ExitIdentifierSeq(c *IdentifierSeqContext)

	// ExitOptScanParams is called when exiting the optScanParams production.
	ExitOptScanParams(c *OptScanParamsContext)

	// ExitTableName is called when exiting the tableName production.
	ExitTableName(c *TableNameContext)

	// ExitAliasedQuery is called when exiting the aliasedQuery production.
	ExitAliasedQuery(c *AliasedQueryContext)

	// ExitTableValuedFunction is called when exiting the tableValuedFunction production.
	ExitTableValuedFunction(c *TableValuedFunctionContext)

	// ExitRelationList is called when exiting the relationList production.
	ExitRelationList(c *RelationListContext)

	// ExitMaterializedViewName is called when exiting the materializedViewName production.
	ExitMaterializedViewName(c *MaterializedViewNameContext)

	// ExitPropertyClause is called when exiting the propertyClause production.
	ExitPropertyClause(c *PropertyClauseContext)

	// ExitPropertyItemList is called when exiting the propertyItemList production.
	ExitPropertyItemList(c *PropertyItemListContext)

	// ExitPropertyItem is called when exiting the propertyItem production.
	ExitPropertyItem(c *PropertyItemContext)

	// ExitPropertyKey is called when exiting the propertyKey production.
	ExitPropertyKey(c *PropertyKeyContext)

	// ExitPropertyValue is called when exiting the propertyValue production.
	ExitPropertyValue(c *PropertyValueContext)

	// ExitTableAlias is called when exiting the tableAlias production.
	ExitTableAlias(c *TableAliasContext)

	// ExitMultipartIdentifier is called when exiting the multipartIdentifier production.
	ExitMultipartIdentifier(c *MultipartIdentifierContext)

	// ExitSimpleColumnDefs is called when exiting the simpleColumnDefs production.
	ExitSimpleColumnDefs(c *SimpleColumnDefsContext)

	// ExitSimpleColumnDef is called when exiting the simpleColumnDef production.
	ExitSimpleColumnDef(c *SimpleColumnDefContext)

	// ExitColumnDefs is called when exiting the columnDefs production.
	ExitColumnDefs(c *ColumnDefsContext)

	// ExitColumnDef is called when exiting the columnDef production.
	ExitColumnDef(c *ColumnDefContext)

	// ExitIndexDefs is called when exiting the indexDefs production.
	ExitIndexDefs(c *IndexDefsContext)

	// ExitIndexDef is called when exiting the indexDef production.
	ExitIndexDef(c *IndexDefContext)

	// ExitPartitionsDef is called when exiting the partitionsDef production.
	ExitPartitionsDef(c *PartitionsDefContext)

	// ExitPartitionDef is called when exiting the partitionDef production.
	ExitPartitionDef(c *PartitionDefContext)

	// ExitLessThanPartitionDef is called when exiting the lessThanPartitionDef production.
	ExitLessThanPartitionDef(c *LessThanPartitionDefContext)

	// ExitFixedPartitionDef is called when exiting the fixedPartitionDef production.
	ExitFixedPartitionDef(c *FixedPartitionDefContext)

	// ExitStepPartitionDef is called when exiting the stepPartitionDef production.
	ExitStepPartitionDef(c *StepPartitionDefContext)

	// ExitInPartitionDef is called when exiting the inPartitionDef production.
	ExitInPartitionDef(c *InPartitionDefContext)

	// ExitPartitionValueList is called when exiting the partitionValueList production.
	ExitPartitionValueList(c *PartitionValueListContext)

	// ExitPartitionValueDef is called when exiting the partitionValueDef production.
	ExitPartitionValueDef(c *PartitionValueDefContext)

	// ExitRollupDefs is called when exiting the rollupDefs production.
	ExitRollupDefs(c *RollupDefsContext)

	// ExitRollupDef is called when exiting the rollupDef production.
	ExitRollupDef(c *RollupDefContext)

	// ExitAggTypeDef is called when exiting the aggTypeDef production.
	ExitAggTypeDef(c *AggTypeDefContext)

	// ExitTabletList is called when exiting the tabletList production.
	ExitTabletList(c *TabletListContext)

	// ExitInlineTable is called when exiting the inlineTable production.
	ExitInlineTable(c *InlineTableContext)

	// ExitNamedExpression is called when exiting the namedExpression production.
	ExitNamedExpression(c *NamedExpressionContext)

	// ExitNamedExpressionSeq is called when exiting the namedExpressionSeq production.
	ExitNamedExpressionSeq(c *NamedExpressionSeqContext)

	// ExitExpression is called when exiting the expression production.
	ExitExpression(c *ExpressionContext)

	// ExitLambdaExpression is called when exiting the lambdaExpression production.
	ExitLambdaExpression(c *LambdaExpressionContext)

	// ExitExist is called when exiting the exist production.
	ExitExist(c *ExistContext)

	// ExitLogicalNot is called when exiting the logicalNot production.
	ExitLogicalNot(c *LogicalNotContext)

	// ExitPredicated is called when exiting the predicated production.
	ExitPredicated(c *PredicatedContext)

	// ExitIsnull is called when exiting the isnull production.
	ExitIsnull(c *IsnullContext)

	// ExitIs_not_null_pred is called when exiting the is_not_null_pred production.
	ExitIs_not_null_pred(c *Is_not_null_predContext)

	// ExitLogicalBinary is called when exiting the logicalBinary production.
	ExitLogicalBinary(c *LogicalBinaryContext)

	// ExitDoublePipes is called when exiting the doublePipes production.
	ExitDoublePipes(c *DoublePipesContext)

	// ExitRowConstructor is called when exiting the rowConstructor production.
	ExitRowConstructor(c *RowConstructorContext)

	// ExitRowConstructorItem is called when exiting the rowConstructorItem production.
	ExitRowConstructorItem(c *RowConstructorItemContext)

	// ExitPredicate is called when exiting the predicate production.
	ExitPredicate(c *PredicateContext)

	// ExitValueExpressionDefault is called when exiting the valueExpressionDefault production.
	ExitValueExpressionDefault(c *ValueExpressionDefaultContext)

	// ExitComparison is called when exiting the comparison production.
	ExitComparison(c *ComparisonContext)

	// ExitBitOperation is called when exiting the bitOperation production.
	ExitBitOperation(c *BitOperationContext)

	// ExitArithmeticBinary is called when exiting the arithmeticBinary production.
	ExitArithmeticBinary(c *ArithmeticBinaryContext)

	// ExitArithmeticUnary is called when exiting the arithmeticUnary production.
	ExitArithmeticUnary(c *ArithmeticUnaryContext)

	// ExitDatetimeUnit is called when exiting the datetimeUnit production.
	ExitDatetimeUnit(c *DatetimeUnitContext)

	// ExitDateCeil is called when exiting the dateCeil production.
	ExitDateCeil(c *DateCeilContext)

	// ExitDereference is called when exiting the dereference production.
	ExitDereference(c *DereferenceContext)

	// ExitCurrentDate is called when exiting the currentDate production.
	ExitCurrentDate(c *CurrentDateContext)

	// ExitTimestampadd is called when exiting the timestampadd production.
	ExitTimestampadd(c *TimestampaddContext)

	// ExitDate_sub is called when exiting the date_sub production.
	ExitDate_sub(c *Date_subContext)

	// ExitCast is called when exiting the cast production.
	ExitCast(c *CastContext)

	// ExitParenthesizedExpression is called when exiting the parenthesizedExpression production.
	ExitParenthesizedExpression(c *ParenthesizedExpressionContext)

	// ExitUserVariable is called when exiting the userVariable production.
	ExitUserVariable(c *UserVariableContext)

	// ExitElementAt is called when exiting the elementAt production.
	ExitElementAt(c *ElementAtContext)

	// ExitLocalTimestamp is called when exiting the localTimestamp production.
	ExitLocalTimestamp(c *LocalTimestampContext)

	// ExitCharFunction is called when exiting the charFunction production.
	ExitCharFunction(c *CharFunctionContext)

	// ExitIntervalLiteral is called when exiting the intervalLiteral production.
	ExitIntervalLiteral(c *IntervalLiteralContext)

	// ExitArrayRange is called when exiting the arrayRange production.
	ExitArrayRange(c *ArrayRangeContext)

	// ExitSimpleCase is called when exiting the simpleCase production.
	ExitSimpleCase(c *SimpleCaseContext)

	// ExitColumnReference is called when exiting the columnReference production.
	ExitColumnReference(c *ColumnReferenceContext)

	// ExitStar is called when exiting the star production.
	ExitStar(c *StarContext)

	// ExitConvertType is called when exiting the convertType production.
	ExitConvertType(c *ConvertTypeContext)

	// ExitTimestampdiff is called when exiting the timestampdiff production.
	ExitTimestampdiff(c *TimestampdiffContext)

	// ExitConvertCharSet is called when exiting the convertCharSet production.
	ExitConvertCharSet(c *ConvertCharSetContext)

	// ExitSubqueryExpression is called when exiting the subqueryExpression production.
	ExitSubqueryExpression(c *SubqueryExpressionContext)

	// ExitEncryptKey is called when exiting the encryptKey production.
	ExitEncryptKey(c *EncryptKeyContext)

	// ExitDate_add is called when exiting the date_add production.
	ExitDate_add(c *Date_addContext)

	// ExitCurrentTime is called when exiting the currentTime production.
	ExitCurrentTime(c *CurrentTimeContext)

	// ExitLocalTime is called when exiting the localTime production.
	ExitLocalTime(c *LocalTimeContext)

	// ExitSystemVariable is called when exiting the systemVariable production.
	ExitSystemVariable(c *SystemVariableContext)

	// ExitCollate is called when exiting the collate production.
	ExitCollate(c *CollateContext)

	// ExitCurrentUser is called when exiting the currentUser production.
	ExitCurrentUser(c *CurrentUserContext)

	// ExitConstantDefault is called when exiting the constantDefault production.
	ExitConstantDefault(c *ConstantDefaultContext)

	// ExitExtract is called when exiting the extract production.
	ExitExtract(c *ExtractContext)

	// ExitCurrentTimestamp is called when exiting the currentTimestamp production.
	ExitCurrentTimestamp(c *CurrentTimestampContext)

	// ExitFunctionCall is called when exiting the functionCall production.
	ExitFunctionCall(c *FunctionCallContext)

	// ExitArraySlice is called when exiting the arraySlice production.
	ExitArraySlice(c *ArraySliceContext)

	// ExitDateFloor is called when exiting the dateFloor production.
	ExitDateFloor(c *DateFloorContext)

	// ExitSearchedCase is called when exiting the searchedCase production.
	ExitSearchedCase(c *SearchedCaseContext)

	// ExitExcept is called when exiting the except production.
	ExitExcept(c *ExceptContext)

	// ExitReplace is called when exiting the replace production.
	ExitReplace(c *ReplaceContext)

	// ExitCastDataType is called when exiting the castDataType production.
	ExitCastDataType(c *CastDataTypeContext)

	// ExitFunctionCallExpression is called when exiting the functionCallExpression production.
	ExitFunctionCallExpression(c *FunctionCallExpressionContext)

	// ExitFunctionIdentifier is called when exiting the functionIdentifier production.
	ExitFunctionIdentifier(c *FunctionIdentifierContext)

	// ExitFunctionNameIdentifier is called when exiting the functionNameIdentifier production.
	ExitFunctionNameIdentifier(c *FunctionNameIdentifierContext)

	// ExitWindowSpec is called when exiting the windowSpec production.
	ExitWindowSpec(c *WindowSpecContext)

	// ExitWindowFrame is called when exiting the windowFrame production.
	ExitWindowFrame(c *WindowFrameContext)

	// ExitFrameUnits is called when exiting the frameUnits production.
	ExitFrameUnits(c *FrameUnitsContext)

	// ExitFrameBoundary is called when exiting the frameBoundary production.
	ExitFrameBoundary(c *FrameBoundaryContext)

	// ExitQualifiedName is called when exiting the qualifiedName production.
	ExitQualifiedName(c *QualifiedNameContext)

	// ExitSpecifiedPartition is called when exiting the specifiedPartition production.
	ExitSpecifiedPartition(c *SpecifiedPartitionContext)

	// ExitNullLiteral is called when exiting the nullLiteral production.
	ExitNullLiteral(c *NullLiteralContext)

	// ExitTypeConstructor is called when exiting the typeConstructor production.
	ExitTypeConstructor(c *TypeConstructorContext)

	// ExitNumericLiteral is called when exiting the numericLiteral production.
	ExitNumericLiteral(c *NumericLiteralContext)

	// ExitBooleanLiteral is called when exiting the booleanLiteral production.
	ExitBooleanLiteral(c *BooleanLiteralContext)

	// ExitStringLiteral is called when exiting the stringLiteral production.
	ExitStringLiteral(c *StringLiteralContext)

	// ExitArrayLiteral is called when exiting the arrayLiteral production.
	ExitArrayLiteral(c *ArrayLiteralContext)

	// ExitMapLiteral is called when exiting the mapLiteral production.
	ExitMapLiteral(c *MapLiteralContext)

	// ExitStructLiteral is called when exiting the structLiteral production.
	ExitStructLiteral(c *StructLiteralContext)

	// ExitPlaceholder is called when exiting the placeholder production.
	ExitPlaceholder(c *PlaceholderContext)

	// ExitComparisonOperator is called when exiting the comparisonOperator production.
	ExitComparisonOperator(c *ComparisonOperatorContext)

	// ExitBooleanValue is called when exiting the booleanValue production.
	ExitBooleanValue(c *BooleanValueContext)

	// ExitWhenClause is called when exiting the whenClause production.
	ExitWhenClause(c *WhenClauseContext)

	// ExitInterval is called when exiting the interval production.
	ExitInterval(c *IntervalContext)

	// ExitUnitIdentifier is called when exiting the unitIdentifier production.
	ExitUnitIdentifier(c *UnitIdentifierContext)

	// ExitDataTypeWithNullable is called when exiting the dataTypeWithNullable production.
	ExitDataTypeWithNullable(c *DataTypeWithNullableContext)

	// ExitComplexDataType is called when exiting the complexDataType production.
	ExitComplexDataType(c *ComplexDataTypeContext)

	// ExitAggStateDataType is called when exiting the aggStateDataType production.
	ExitAggStateDataType(c *AggStateDataTypeContext)

	// ExitPrimitiveDataType is called when exiting the primitiveDataType production.
	ExitPrimitiveDataType(c *PrimitiveDataTypeContext)

	// ExitPrimitiveColType is called when exiting the primitiveColType production.
	ExitPrimitiveColType(c *PrimitiveColTypeContext)

	// ExitComplexColTypeList is called when exiting the complexColTypeList production.
	ExitComplexColTypeList(c *ComplexColTypeListContext)

	// ExitComplexColType is called when exiting the complexColType production.
	ExitComplexColType(c *ComplexColTypeContext)

	// ExitCommentSpec is called when exiting the commentSpec production.
	ExitCommentSpec(c *CommentSpecContext)

	// ExitSample is called when exiting the sample production.
	ExitSample(c *SampleContext)

	// ExitSampleByPercentile is called when exiting the sampleByPercentile production.
	ExitSampleByPercentile(c *SampleByPercentileContext)

	// ExitSampleByRows is called when exiting the sampleByRows production.
	ExitSampleByRows(c *SampleByRowsContext)

	// ExitTableSnapshot is called when exiting the tableSnapshot production.
	ExitTableSnapshot(c *TableSnapshotContext)

	// ExitErrorCapturingIdentifier is called when exiting the errorCapturingIdentifier production.
	ExitErrorCapturingIdentifier(c *ErrorCapturingIdentifierContext)

	// ExitErrorIdent is called when exiting the errorIdent production.
	ExitErrorIdent(c *ErrorIdentContext)

	// ExitRealIdent is called when exiting the realIdent production.
	ExitRealIdent(c *RealIdentContext)

	// ExitIdentifier is called when exiting the identifier production.
	ExitIdentifier(c *IdentifierContext)

	// ExitUnquotedIdentifier is called when exiting the unquotedIdentifier production.
	ExitUnquotedIdentifier(c *UnquotedIdentifierContext)

	// ExitQuotedIdentifierAlternative is called when exiting the quotedIdentifierAlternative production.
	ExitQuotedIdentifierAlternative(c *QuotedIdentifierAlternativeContext)

	// ExitQuotedIdentifier is called when exiting the quotedIdentifier production.
	ExitQuotedIdentifier(c *QuotedIdentifierContext)

	// ExitIntegerLiteral is called when exiting the integerLiteral production.
	ExitIntegerLiteral(c *IntegerLiteralContext)

	// ExitDecimalLiteral is called when exiting the decimalLiteral production.
	ExitDecimalLiteral(c *DecimalLiteralContext)

	// ExitNonReserved is called when exiting the nonReserved production.
	ExitNonReserved(c *NonReservedContext)
}
