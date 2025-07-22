#!/usr/bin/env python3
"""
Python Test Runner for URL-DB MCP Server
Manages and executes all Python integration tests for the URL-DB project.
"""

import os
import sys
import subprocess
import json
import time
import argparse
from pathlib import Path
from typing import Dict, List, Optional, Tuple
import logging

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

class TestRunner:
    """Manages execution of Python integration tests for URL-DB MCP server."""
    
    def __init__(self, server_binary: str = "../../bin/url-db", db_path: str = ":memory:"):
        self.server_binary = Path(__file__).parent / server_binary
        self.db_path = db_path
        self.test_results: Dict[str, Dict] = {}
        
        # Test categories and their files
        self.test_categories = {
            "demo": [
                "server_info_demo.py",
                "list_domains.py"
            ],
            "basic": [
                "test_mcp_client.py",
                "test_node_attributes.py"
            ],
            "comprehensive": [
                "test_all_mcp_tools.py",
                "test_mcp_tools.py",
                "test_mcp_final.py"
            ],
            "domain_attributes": [
                "test_domain_attributes.py",
                "test_mcp_domain_attributes.py",
                "test_final.py",
                "test_mcp_persistent.py"
            ],
            "advanced": [
                "test_scenarios.py"
            ]
        }
    
    def check_server_binary(self) -> bool:
        """Check if the server binary exists."""
        if not self.server_binary.exists():
            logger.error(f"Server binary not found: {self.server_binary}")
            logger.info("Run 'make build' to build the server first")
            return False
        return True
    
    def start_server(self) -> Optional[subprocess.Popen]:
        """Start the MCP server in stdio mode."""
        try:
            cmd = [str(self.server_binary), "-mcp-mode=stdio", f"-db-path={self.db_path}"]
            logger.info(f"Starting server: {' '.join(cmd)}")
            
            process = subprocess.Popen(
                cmd,
                stdin=subprocess.PIPE,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=True
            )
            
            # Give server time to start
            time.sleep(1)
            
            # Check if process is still running
            if process.poll() is not None:
                stderr = process.stderr.read() if process.stderr else ""
                logger.error(f"Server failed to start. Error: {stderr}")
                return None
                
            logger.info("Server started successfully")
            return process
            
        except Exception as e:
            logger.error(f"Failed to start server: {e}")
            return None
    
    def run_test(self, test_file: str) -> Tuple[bool, str, float]:
        """Run a single test file."""
        start_time = time.time()
        test_path = Path(__file__).parent / test_file
        
        if not test_path.exists():
            return False, f"Test file not found: {test_file}", 0.0
        
        try:
            logger.info(f"Running test: {test_file}")
            result = subprocess.run(
                [sys.executable, str(test_path)],
                capture_output=True,
                text=True,
                timeout=30,
                cwd=Path(__file__).parent
            )
            
            duration = time.time() - start_time
            
            if result.returncode == 0:
                logger.info(f"✅ {test_file} passed ({duration:.2f}s)")
                return True, result.stdout, duration
            else:
                logger.error(f"❌ {test_file} failed ({duration:.2f}s)")
                return False, f"STDOUT:\n{result.stdout}\nSTDERR:\n{result.stderr}", duration
                
        except subprocess.TimeoutExpired:
            duration = time.time() - start_time
            logger.error(f"❌ {test_file} timed out ({duration:.2f}s)")
            return False, "Test timed out after 30 seconds", duration
        except Exception as e:
            duration = time.time() - start_time
            logger.error(f"❌ {test_file} error ({duration:.2f}s): {e}")
            return False, str(e), duration
    
    def run_category(self, category: str) -> Dict:
        """Run all tests in a category."""
        if category not in self.test_categories:
            logger.error(f"Unknown category: {category}")
            return {"success": False, "error": f"Unknown category: {category}"}
        
        logger.info(f"\n🏃 Running {category} tests...")
        results = {
            "category": category,
            "tests": {},
            "passed": 0,
            "failed": 0,
            "total_duration": 0.0
        }
        
        for test_file in self.test_categories[category]:
            success, output, duration = self.run_test(test_file)
            
            results["tests"][test_file] = {
                "success": success,
                "output": output,
                "duration": duration
            }
            results["total_duration"] += duration
            
            if success:
                results["passed"] += 1
            else:
                results["failed"] += 1
        
        return results
    
    def run_all_tests(self, categories: Optional[List[str]] = None) -> Dict:
        """Run tests for specified categories or all categories."""
        if categories is None:
            categories = list(self.test_categories.keys())
        
        logger.info("🚀 Starting Python integration tests for URL-DB MCP server")
        logger.info("=" * 60)
        
        # Check server binary
        if not self.check_server_binary():
            return {"success": False, "error": "Server binary not found"}
        
        overall_results = {
            "success": True,
            "categories": {},
            "summary": {
                "total_tests": 0,
                "passed": 0,
                "failed": 0,
                "total_duration": 0.0
            }
        }
        
        for category in categories:
            if category not in self.test_categories:
                logger.warning(f"Skipping unknown category: {category}")
                continue
                
            category_results = self.run_category(category)
            overall_results["categories"][category] = category_results
            
            # Update summary
            overall_results["summary"]["total_tests"] += len(self.test_categories[category])
            overall_results["summary"]["passed"] += category_results["passed"]
            overall_results["summary"]["failed"] += category_results["failed"]
            overall_results["summary"]["total_duration"] += category_results["total_duration"]
            
            if category_results["failed"] > 0:
                overall_results["success"] = False
        
        return overall_results
    
    def print_summary(self, results: Dict):
        """Print test results summary."""
        summary = results["summary"]
        
        print("\n" + "=" * 60)
        print("📊 TEST RESULTS SUMMARY")
        print("=" * 60)
        
        print(f"Total Tests: {summary['total_tests']}")
        print(f"✅ Passed: {summary['passed']}")
        print(f"❌ Failed: {summary['failed']}")
        print(f"⏱️  Total Duration: {summary['total_duration']:.2f}s")
        print(f"📈 Success Rate: {(summary['passed'] / summary['total_tests'] * 100):.1f}%")
        
        if results["success"]:
            print("\n🎉 All tests passed!")
        else:
            print("\n💥 Some tests failed. Check logs above for details.")
            
        print("\n📝 Category Breakdown:")
        for category, cat_results in results["categories"].items():
            status = "✅" if cat_results["failed"] == 0 else "❌"
            print(f"  {status} {category}: {cat_results['passed']}/{cat_results['passed'] + cat_results['failed']} passed")
    
    def list_tests(self):
        """List all available tests by category."""
        print("📋 Available Tests by Category:")
        print("=" * 40)
        
        for category, tests in self.test_categories.items():
            print(f"\n🏷️  {category.upper()}:")
            for test in tests:
                print(f"   • {test}")
        
        print(f"\nTotal: {sum(len(tests) for tests in self.test_categories.values())} tests")

def main():
    parser = argparse.ArgumentParser(description="Run Python integration tests for URL-DB MCP server")
    parser.add_argument("--category", "-c", choices=["demo", "basic", "comprehensive", "domain_attributes", "advanced", "all"], 
                       default="all", help="Test category to run")
    parser.add_argument("--list", "-l", action="store_true", help="List all available tests")
    parser.add_argument("--server-binary", "-s", default="../../bin/url-db", 
                       help="Path to server binary")
    parser.add_argument("--db-path", "-d", default=":memory:", 
                       help="Database path (use :memory: for in-memory)")
    parser.add_argument("--verbose", "-v", action="store_true", help="Verbose output")
    
    args = parser.parse_args()
    
    if args.verbose:
        logging.getLogger().setLevel(logging.DEBUG)
    
    runner = TestRunner(args.server_binary, args.db_path)
    
    if args.list:
        runner.list_tests()
        return
    
    # Run tests
    if args.category == "all":
        categories = None
    else:
        categories = [args.category]
    
    results = runner.run_all_tests(categories)
    runner.print_summary(results)
    
    # Exit with appropriate code
    sys.exit(0 if results["success"] else 1)

if __name__ == "__main__":
    main()